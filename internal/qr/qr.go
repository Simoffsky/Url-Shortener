package qr

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/qr/cache"

	"url-shorter/pkg/log"
	pb "url-shorter/pkg/proto/qr"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

type QRServer struct {
	pb.UnimplementedQRCodeServiceServer
	service QRService
	config  config.Config
	logger  log.Logger
}

func NewQRServer(config config.Config) *QRServer {
	return &QRServer{
		config: config,
		logger: log.NewDefaultLogger(log.LevelFromString(config.LoggerLevel)),
	}
}

func (s *QRServer) Start() error {
	if err := s.configure(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", s.config.QrAddr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterQRCodeServiceServer(grpcServer, s)

	//Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("Starting gRPC server on address: " + s.config.QrAddr)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve gRPC: %s" + err.Error())
		}
	}()

	<-sigChan

	s.logger.Info("Shutting down gRPC server gracefully")
	stopCh := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopCh)
	}()

	select {
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
		s.logger.Info("gRPC server stopped by timeout")
	case <-stopCh:
		s.logger.Info("gRPC server stopped gracefully")
	}
	return nil
}

func (s *QRServer) configure() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     s.config.RedisAddr,
		Password: "",
		DB:       0,
	})

	err := redisClient.Ping().Err()

	if err != nil {
		return err
	}

	s.logger.Info("Connected to Redis")

	s.service = NewDefaultQRService(cache.NewRedisCache(redisClient), s.logger)

	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel)).
		WithTimePrefix(time.DateTime)

	return nil
}
