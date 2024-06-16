package main

import (
	"log"
	"url-shorter/internal/config"
	"url-shorter/internal/server"
)

func main() {
	config := config.NewEnvConfig()

	app := server.NewLinkServer(config)

	err := app.Start()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
