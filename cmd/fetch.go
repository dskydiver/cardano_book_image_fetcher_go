/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"cardano-book-image-fetcher/pkg"
	"cardano-book-image-fetcher/pkg/api"

	"github.com/blockfrost/blockfrost-go"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		policyID, err := cmd.Flags().GetString("policy-id")
		if err != nil {
			fmt.Println("The policy-id flag is invalid")
			return
		}
		outputDir, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			fmt.Println("The output-dir flag is invalid")
			return
		}
		fmt.Println("Policy ID:", policyID)
		fmt.Println("Output directory:", outputDir)
		bookioClient := &api.BookIOClient{
			BaseURL: "https://api.book.io",
		}
		collections, err := bookioClient.FetchCollections()
		if err != nil {
			fmt.Println("Error fetching collections:", err)
			return
		}
		is_valid := pkg.VerifyPolicyID(policyID, collections)
		if !is_valid {
			fmt.Println("Policy ID is not valid")
			return
		} else {
			fmt.Println("Policy ID is valid")
		}
		api := blockfrost.NewAPIClient(blockfrost.APIClientOptions{})
		err = pkg.FetchImages(api, policyID, outputDir)
		if err != nil {
			fmt.Println("Error fetching images:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.PersistentFlags().String("policy-id", "", "Policy ID")
	fetchCmd.PersistentFlags().String("output-dir", "", "Output directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
