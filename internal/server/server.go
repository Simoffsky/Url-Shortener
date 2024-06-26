package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/kafka"
	repository "url-shorter/internal/repository/links"
	"url-shorter/internal/server/services"
	"url-shorter/pkg/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LinkServer struct {
	config      config.Config
	linkService services.LinkService
	authService services.AuthService
	statService services.StatService
	logger      log.Logger
}

func NewLinkServer(config config.Config) *LinkServer {
	return &LinkServer{
		config: config,
	}
}

func (s *LinkServer) configureServer() error {
	var err error

	pgConn, err := repository.ConnectToDB(s.config.DbConn)
	if err != nil {
		return err
	}
	linksRepo := repository.NewPgLinksRepository(pgConn)
	qrRepo, err := repository.NewQrGRPCRepository(s.config.QrAddr)
	if err != nil {
		return err
	}

	s.linkService = services.NewDefaultLinkService(linksRepo, qrRepo)
	authService, err := services.NewAuthServiceGRPC(s.config.AuthAddr, s.config.JwtSecret)
	if err != nil {
		return err
	}
	s.authService = authService
	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel),
	).WithTimePrefix(time.DateTime)

	kafkaManager, err := kafka.NewProducerManager(s.config.KafkaBrokers, s.config.KafkaTopic, s.logger)
	if err != nil {
		return err
	}
	s.statService, err = services.NewStatServiceGRPC(s.config.StatsAddr, &kafkaManager)

	kafkaManager.RunProducer(context.Background())
	if err != nil {
		return err
	}
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

	mux.HandleFunc("/create-url/", s.handleCreateLink)
	mux.Handle("/{short}/", WithMetrics(http.HandlerFunc(s.handleLink)))
	mux.Handle("/{short}", WithMetrics(http.HandlerFunc(s.handleLink)))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/qr/", s.handleQRCode)
	mux.HandleFunc("/stats/{short}", s.handleStats)
	mux.HandleFunc("/login/", s.handleLogin)
	mux.HandleFunc("/register/", s.handleRegister)

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
