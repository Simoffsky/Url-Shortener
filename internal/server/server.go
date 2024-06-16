package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/repository"
	"url-shorter/pkg/log"
)

type LinkServer struct {
	config   config.Config
	linkRepo repository.LinksRepository
	logger   log.Logger
}

func NewLinkServer(config config.Config) *LinkServer {
	return &LinkServer{
		config: config,
	}
}

func (s *LinkServer) configureServer() error {
	s.linkRepo = repository.NewMemoryLinksRepository()
	s.logger = log.NewDefaultLogger(log.LevelFromString(s.config.LoggerLevel))
	return nil
}

func (s *LinkServer) Start() error {
	err := s.configureServer()
	if err != nil {
		return err
	}

	return s.startHTTPServer()
}

func (s *LinkServer) startHTTPServer() error {
	http.HandleFunc("POST /create-url/", handler)
	config := config.NewEnvConfig()
	s.logger.Debug("Config parameters: " + config.String())

	server := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        nil,
		ReadTimeout:    config.HTTPTimeout,
		WriteTimeout:   config.HTTPTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//Graceful shutdown
	errCh := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		s.logger.Info("Server started on port " + s.config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-quit:
		s.logger.Info("Server is shutting down...")
	case err := <-errCh:
		s.logger.Error("HTTP Server error:" + err.Error())
	}

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return server.Shutdown(ctx)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := log.NewDefaultLogger(log.Info)
	data := []byte("Hello, World!")
	n, err := w.Write(data)
	if err != nil {
		logger.Error("writing response: " + err.Error())
	}
	fmt.Printf("Wrote %d bytes, expected %d\n", n, len(data))
}
