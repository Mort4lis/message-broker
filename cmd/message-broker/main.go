package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Mort4lis/message-broker/internal/app"
)

func main() {
	var confPath string

	flag.StringVar(&confPath, "c", "./configs/config.yaml", "The configuration file path")
	flag.Parse()

	if err := app.Run(confPath); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurs while running the application: %v", err)
		os.Exit(1)
	}
}
