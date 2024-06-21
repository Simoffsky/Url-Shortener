package services

import (
	"context"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(models.User) error
	Login(models.User) (string, error)
}

type AuthServiceGRPC struct {
	authClient pb.AuthServiceClient
	jwtSecret  string
}

func NewAuthServiceGRPC(authAddr string, jwtSecret string) (*AuthServiceGRPC, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(authAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthServiceGRPC{
		authClient: client,
		jwtSecret:  jwtSecret,
	}, nil
}

func (s *AuthServiceGRPC) Register(user models.User) error {
	_, err := s.authClient.Register(context.Background(), &pb.RegisterRequest{
		Login:    user.Login,
		Password: user.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			return models.ErrUserExists
		}
	}
	return err
}

// returns token
func (s *AuthServiceGRPC) Login(user models.User) (string, error) {
	resp, err := s.authClient.Login(context.Background(), &pb.LoginRequest{
		Login:    user.Login,
		Password: user.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			return "", models.ErrInvalidPassword
		}
		if ok && st.Code() == codes.NotFound {
			return "", models.ErrUserNotFound
		}
		return "", err
	}

	return resp.Token, nil
}
