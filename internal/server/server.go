package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/repository"
	"url-shorter/pkg/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LinkServer struct {
	config      config.Config
	linkService LinkService
	logger      log.Logger
}

func NewLinkServer(config config.Config) *LinkServer {
	return &LinkServer{
		config: config,
	}
}

func (s *LinkServer) configureServer() error {
	var err error
	s.linkService, err = NewDefaultLinkService(repository.NewMemoryLinksRepository(), ":"+s.config.QRGRPCPort)
	if err != nil {
		return err
	}
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
	mux := http.NewServeMux()

	mux.HandleFunc("/create-url/", s.handler)
	mux.Handle("/", WithMetrics(http.HandlerFunc(s.handleRedirect)))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/qr/", s.handleQRCode)

	s.logger.Debug("Config parameters: " + s.config.String())

	server := &http.Server{
		Addr:           ":" + s.config.ServerPort,
		Handler:        mux,
		ReadTimeout:    s.config.HTTPTimeout,
		WriteTimeout:   s.config.HTTPTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//Graceful shutdown
	errCh := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		s.logger.Info("Server started on port " + s.config.ServerPort)
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
