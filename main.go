/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cardano-book-image-fetcher/cmd"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	cmd.Execute()
}
