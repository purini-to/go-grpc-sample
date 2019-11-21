package main

import (
	"log"
	"os"

	"github.com/purini-to/go-grpc-sample/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Printf("Faild command. error: %v", err)
		os.Exit(1)
	}
}
