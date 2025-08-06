package main

import (
	"fmt"
	"os"

	"github.com/YusovID/file-downloader-06.08.2025/internal/config"
)

func init() {
	err := os.Setenv("CONFIG_PATH", "./config/local.yaml")
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger

	// TODO: init router

	// TODO: init server
}
