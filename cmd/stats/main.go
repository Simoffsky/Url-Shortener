package main

import (
	"log"
	"url-shorter/internal/config"
	"url-shorter/internal/stats"
)

func main() {
	config := config.NewEnvConfig()

	app := stats.NewStatsServer(config)

	err := app.Start()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
