/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// timezoneCmd represents the timezone command
var timezoneCmd = &cobra.Command{
	Use:   "timezone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		timezone := args[0]
		location, err := time.LoadLocation(timezone)
		if err != nil {
			log.Fatalln("The timezone string is invalid")
		}
		dateFlag, err := cmd.Flags().GetString("date")
		if err != nil {
			log.Fatalln("The date flag is invalid")
		}
		var date string
		if dateFlag != "" {
			date = time.Now().In(location).Format(dateFlag)
		} else {
			date = time.Now().In(location).Format(time.RFC1123)[:10]
		}
		fmt.Printf("Current date in %v: %v\n", timezone, date)
	},
}

func init() {
	rootCmd.AddCommand(timezoneCmd)
	timezoneCmd.PersistentFlags().String("date", "", "returns the date in a time zone in a specified format")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timezoneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timezoneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
