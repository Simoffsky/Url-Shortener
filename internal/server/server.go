package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/models"
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
	http.HandleFunc("POST /create-url/", s.handler)
	http.HandleFunc("GET /", s.handleRedirect)
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

type Request struct {
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
}

type Response struct {
	ShortUrl string `json:"short_url"`
}

// TODO: Check url format
func (s *LinkServer) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Url == "" {
		http.Error(w, "Url is required", http.StatusBadRequest)
		return
	}

	if req.ShortUrl == "" {
		http.Error(w, "ShortUrl is required", http.StatusBadRequest)
		return
	}

	err = s.linkRepo.CreateLink(req.Url, req.ShortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Response{
		ShortUrl: req.ShortUrl,
	}
	
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *LinkServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	short := r.URL.Path[len("/"):]

	url, err := s.linkRepo.GetLink(short)
	if err != nil {
		if errors.Is(err, models.ErrLinkNotFound) {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
