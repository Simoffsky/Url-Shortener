package services

import (
	"context"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService interface {
	Register(models.User) error
	Login(models.User) (string, error)
}

type AuthServiceGRPC struct {
	authClient pb.AuthServiceClient
}

func NewAuthServiceGRPC(authAddr string) (*AuthServiceGRPC, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(authAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceGRPC{
		authClient: client,
	}, nil
}

func (s *AuthServiceGRPC) Register(user models.User) error {
	_, err := s.authClient.Register(context.Background(), &pb.RegisterRequest{
		Login:    user.Login,
		Password: user.Password,
	})
	return err
}

func (s *AuthServiceGRPC) Login(user models.User) (string, error) {
	resp, err := s.authClient.Login(context.Background(), &pb.LoginRequest{
		Login:    user.Login,
		Password: user.Password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}
