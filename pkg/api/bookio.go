package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Collection struct {
	CollectionID string `json:"collection_id"`
	Description  string `json:"description"`
	Blockchain   string `json:"blockchain"`
	Network      string `json:"network"`
}

type CollectionResponse struct {
	Type string       `json:"type"`
	Data []Collection `json:"data"`
}

type BookIOClient struct {
	BaseURL string
}

func (c *BookIOClient) FetchCollections() ([]Collection, error) {
	endpoint := fmt.Sprintf("%s/api/v0/collections", c.BaseURL)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyByptes, _ := io.ReadAll(resp.Body)

	var collectionsResp CollectionResponse
	err = json.Unmarshal(bodyByptes, &collectionsResp)
	if err != nil {
		return nil, err
	}

	return collectionsResp.Data, nil
}
