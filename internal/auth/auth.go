package auth

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shorter/internal/config"
	"url-shorter/pkg/log"
	pb "url-shorter/pkg/proto/auth"

	"google.golang.org/grpc"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer

	service AuthService
	config  config.Config
	logger  log.Logger
}

func NewAuthServer(config config.Config) *AuthServer {
	return &AuthServer{
		config: config,
		logger: log.NewDefaultLogger(log.LevelFromString(config.LoggerLevel)),
	}
}

func (s *AuthServer) Start() error {
	if err := s.configure(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", s.config.AuthAddr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcServer, s)

	//Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("Starting gRPC server on address: " + s.config.AuthAddr)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve gRPC: %s" + err.Error())
		}
	}()

	<-sigChan

	grpcServer.GracefulStop()
	s.logger.Info("Shutting down gRPC server gracefully")
	return nil
}

func (s *AuthServer) configure() error {

	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel)).
		WithTimePrefix(time.DateTime)
	return nil
}
