package main

import (
	"fmt"
	"log"
	"net/http"
	"url-shorter/internal/config"
)

func main() {
	http.HandleFunc("/create-url/", handler)
	fmt.Println("server started")
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
	data := []byte("Hello, World!")
	n, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote %d bytes, expected %d\n", n, len(data))
}
