package auth

import (
	"context"
	"errors"
	"fmt"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.service.Register(in.Login, in.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	s.logger.Debug(fmt.Sprintf("Sending gRPC response to register user: %s", in.Login))
	return &pb.RegisterResponse{}, nil
}

func (s *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(in.Login, in.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, "invalid password")
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		// Обработка других возможных ошибок
		return nil, status.Error(codes.Internal, "internal server error")
	}

	s.logger.Debug(fmt.Sprintf("Sending gRPC response with token for user: %s", in.Login))
	return &pb.LoginResponse{Token: token}, nil
}
