package services

import (
	"context"
	"url-shorter/internal/models"
	"url-shorter/internal/repository"
	pb "url-shorter/pkg/proto/qr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LinkService interface {
	GetLink(short string) (*models.Link, error)
	RemoveLink(short string) error
	CreateLink(models.Link) error
	GetQRCode(short string) ([]byte, error)
}

type DefaultLinkService struct {
	qrClient  pb.QRCodeServiceClient
	linksRepo repository.LinksRepository
}

func NewDefaultLinkService(linksRepo repository.LinksRepository, qrAddr string) (*DefaultLinkService, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(qrAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewQRCodeServiceClient(conn)

	return &DefaultLinkService{
		qrClient:  client,
		linksRepo: linksRepo,
	}, nil

}

func (s *DefaultLinkService) GetLink(short string) (*models.Link, error) {
	return s.linksRepo.GetLink(short)
}

func (s *DefaultLinkService) RemoveLink(short string) error {
	return s.linksRepo.RemoveLink(short)
}

func (s *DefaultLinkService) CreateLink(link models.Link) error {
	if err := link.Validate(); err != nil {
		return err
	}

	return s.linksRepo.CreateLink(link)
}

func (s *DefaultLinkService) GetQRCode(short string) ([]byte, error) {

	link, err := s.GetLink(short)
	if err != nil {
		return nil, err
	}

	resp, err := s.qrClient.GetQRCode(context.Background(), &pb.QRCodeRequest{Url: link.Url})
	if err != nil {
		return nil, err
	}
	return resp.QrCode, nil
}
