package main

import (
	"log"
	"online-shop-2N/pkg/config"
	"online-shop-2N/pkg/di"
)

func main() {

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Error to load the config: ", err)
	}

	server, err := di.InitializeApi(cfg)
	if err != nil {
		log.Fatal("Failed to initialize the api: ", err)
	}

	if server.Start(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
