package auth

import (
	"context"
	pb "url-shorter/pkg/proto/auth"
)

func (s *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.service.Register(in.Login, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{}, nil
}

func (s *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(in.Login, in.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{Token: token}, nil
}
