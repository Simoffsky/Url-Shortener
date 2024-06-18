package qr

import (
	"context"
	"errors"
	"net"
	"time"
	"url-shorter/internal/config"
	"url-shorter/internal/models"
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

	//TODO: graceful shutdown
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterQRCodeServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

func (s *QRServer) configure() error {
	s.service = NewDefaultQRService()
	s.logger = log.NewDefaultLogger(
		log.LevelFromString(s.config.LoggerLevel)).
		WithTimePrefix(time.DateTime)

	return nil
}

func (s *QRServer) GetQRCode(ctx context.Context, in *pb.QRCodeRequest) (*pb.QRCodeResponse, error) {
	qr, err := s.service.GetQRCode(in.Url)
	if errors.Is(err, models.ErrLinkNotFound) {
		qr, err = s.service.CreateQRCode(in.Url, 256)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &pb.QRCodeResponse{QrCode: qr}, nil
}
