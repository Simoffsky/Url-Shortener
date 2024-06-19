package server

import (
	"context"
	"url-shorter/internal/models"
	repository "url-shorter/internal/repository/links"
	pb "url-shorter/pkg/proto/qr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LinkService interface {
	GetLink(short string) (*models.Link, error)
	RemoveLink(short string) error
	CreateLink(models.Link) error
	GetQRCode(short string, imgSize int) ([]byte, error)
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

// zero size means default size(512x512)
func (s *DefaultLinkService) GetQRCode(link string, imgSize int) ([]byte, error) {
	if imgSize == 0 {
		imgSize = 512
	}

	resp, err := s.qrClient.GetQRCode(context.Background(), &pb.QRCodeRequest{Url: link, Size: int32(imgSize)})
	if err != nil {
		return nil, err
	}
	return resp.QrCode, nil
}
