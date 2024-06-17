package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/models"
	"url-shorter/internal/repository"
	"url-shorter/internal/services"
	"url-shorter/pkg/log"
)

type LinkServer struct {
	config      config.Config
	linkService services.LinkService
	logger      log.Logger
}

func NewLinkServer(config config.Config) *LinkServer {
	return &LinkServer{
		config: config,
	}
}

func (s *LinkServer) configureServer() error {

	s.linkService = services.NewDefaultLinkService(repository.NewMemoryLinksRepository())
	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel),
	).WithTimePrefix(time.DateTime)
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

func (s *LinkServer) handler(w http.ResponseWriter, r *http.Request) {
	var link models.Link

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.linkService.CreateLink(link); err != nil {
		if errors.Is(err, models.ErrLinkAlreadyExists) {
			s.writeError(w, http.StatusBadRequest, err)
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *LinkServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	short := removeTrailingSlash(r.URL.Path[len("/"):])

	s.logger.Debug("Redirecting to: " + short)

	url, err := s.linkService.GetLink(short)
	if err != nil {
		if errors.Is(err, models.ErrLinkNotFound) {
			s.writeError(w, http.StatusNotFound, err)
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
}

func (s *LinkServer) writeError(w http.ResponseWriter, errCode int, err error) {
	s.logger.Error(fmt.Sprintf("HTTP error(%d): %s", errCode, err.Error()))
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func removeTrailingSlash(short string) string {
	if short == "" {
		return short
	}
	if short[len(short)-1] == '/' {
		return short[:len(short)-1]
	}
	return short
}
