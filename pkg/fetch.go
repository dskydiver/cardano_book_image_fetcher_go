package pkg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blockfrost/blockfrost-go"
)

type Metadata struct {
	Asset string
	Url   string
}

const MAX_DOWNLOADS = 10
const NUM_PROCESSORS = 5

func FetchAssets(client blockfrost.APIClient, policyID string) ([]blockfrost.Asset, error) {
	assets, err := client.AssetsByPolicy(context.TODO(), policyID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func ProduceAssets(done <-chan struct{}, client blockfrost.APIClient, assets []blockfrost.Asset) <-chan string {
	assetc := make(chan string)
	go func() {
		defer close(assetc)
		for _, asset := range assets {
			select {
			case <-done:
				return
			case assetc <- asset.Asset:
			}
		}
	}()
	return assetc
}

func ProcessAssets(done <-chan struct{}, client blockfrost.APIClient, assets <-chan string) <-chan Metadata {
	out := make(chan Metadata)
	go func() {
		defer close(out)
		for asset := range assets {
			assetData, err := client.Asset(context.TODO(), asset)
			if err != nil {
				continue
			}

			if assetData.OnchainMetadata.Name != "" {
				select {
				case out <- Metadata{
					Asset: asset,
					Url:   assetData.OnchainMetadata.Image,
				}:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

func DownloadAssets(done <-chan struct{}, assets <-chan Metadata, outputDir string) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0o755)
		if err != nil {
			fmt.Println("Error creating output directory", err)
		}
	}
	sem := make(chan struct{}, 10)
	defer close(sem)
	var wg sync.WaitGroup
	successfulDownloads := 0

	go func() {
		wg.Wait()
	}()

	for m := range assets {
		if successfulDownloads >= 10 {
			break
		}
		select {
		case <-done:
		default:
			sem <- struct{}{}
			wg.Add(1)
			go func(m Metadata) {
				defer wg.Done()
				defer func() { <-sem }()
				fmt.Println("Downloading asset ", m.Asset)
				if err := DownloadImage(outputDir, m.Url, m.Asset); err != nil {
					fmt.Printf("Error downloading image for %s. Error: %s", m.Asset, err)
				} else {
					successfulDownloads++
				}

			}(m)
		}
	}
}

func DownloadImage(outputDir, url, outputFileName string) error {
	ipfsHash := strings.TrimPrefix(url, "ipfs://")
	pubURL := fmt.Sprintf("https://ipfs.io/ipfs/%s", ipfsHash)

	outFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.png", outputFileName))
	tempOutFilePath := outFilePath + ".part"
	if _, err := os.Stat(outFilePath); err == nil {
		return nil
	}

	resp, err := http.Get(pubURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	outFile, err := os.Create(tempOutFilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp.Body)

	if err != nil {
		os.Remove(tempOutFilePath)
		return err
	}

	outFile.Close()
	err = os.Rename(tempOutFilePath, outFilePath)
	if err != nil {
		os.Remove(tempOutFilePath)
		return err
	}
	return nil
}

func FetchImages(client blockfrost.APIClient, policyID, outputDir string) error {
	done := make(chan struct{})
	defer close(done)

	assets, err := FetchAssets(client, policyID)
	if err != nil {
		fmt.Println("Error fetching assets:", err)
		return err
	}
	assetc := ProduceAssets(done, client, assets)
	c := make(chan Metadata)
	var wg sync.WaitGroup
	wg.Add(NUM_PROCESSORS)

	for i := 0; i < NUM_PROCESSORS; i++ {
		go func() {
			defer wg.Done()
			for m := range ProcessAssets(done, client, assetc) {
				select {
				case c <- m:
				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	DownloadAssets(done, c, outputDir)
	return nil
}
