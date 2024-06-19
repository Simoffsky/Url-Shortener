package qr

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shorter/internal/config"

	"url-shorter/pkg/log"
	pb "url-shorter/pkg/proto/qr"

	"google.golang.org/grpc"
)

// TODO: logger
type QRServer struct {
	pb.UnimplementedQRCodeServiceServer
	service QRService
	config  config.Config
	logger  log.Logger
}

func NewQRServer(config config.Config) *QRServer {
	return &QRServer{config: config}
}

func (s *QRServer) Start() error {
	if err := s.configure(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":"+s.config.QRGRPCPort)
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
		s.logger.Info("Starting gRPC server on port " + s.config.QRGRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve gRPC: %s" + err.Error())
		}
	}()

	<-sigChan

	grpcServer.GracefulStop()
	s.logger.Info("Shutting down gRPC server gracefully")
	return nil
}

func (s *QRServer) configure() error {
	s.service = NewDefaultQRService()
	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel)).
		WithTimePrefix(time.DateTime)

	return nil
}
