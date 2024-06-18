package main

import (
	"fmt"
	"url-shorter/internal/config"
	"url-shorter/internal/qr"
)

func main() {
	config := config.NewEnvConfig()

	server := qr.NewQRServer(config)
	fmt.Println("Server started...")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
