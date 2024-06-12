package main

import (
	"fmt"
	"net/http"
	"url-shorter/internal/config"
)

func main() {
	http.HandleFunc("/create-url/", handler)
	fmt.Println("Server started")
	config := config.NewEnvConfig()
	fmt.Printf("%+v\n", config)

	server := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        nil,
		ReadTimeout:    config.HTTPTimeout,
		WriteTimeout:   config.HTTPTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
