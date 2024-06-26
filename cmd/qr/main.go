package main

import (
	"url-shorter/internal/config"
	"url-shorter/internal/qr"
)

func main() {
	config := config.NewEnvConfig()

	server := qr.NewQRServer(config)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
