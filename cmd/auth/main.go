package main

import (
	"url-shorter/internal/auth"
	"url-shorter/internal/config"
)

func main() {
	config := config.NewEnvConfig()

	server := auth.NewAuthServer(config)

	err := server.Start()
	if err != nil {
		panic(err)
	}

}
