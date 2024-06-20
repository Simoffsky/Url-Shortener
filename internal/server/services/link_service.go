package services

import (
	"url-shorter/internal/models"
	repository "url-shorter/internal/repository/links"
)

type LinkService interface {
	GetLink(short string) (*models.Link, error)
	RemoveLink(short string) error
	CreateLink(models.Link) error
	GetQRCode(short string, imgSize int) ([]byte, error)
}

type DefaultLinkService struct {
	qrRepository repository.QrRepository
	linksRepo    repository.LinksRepository
}

func NewDefaultLinkService(linksRepo repository.LinksRepository, qrRepo repository.QrRepository) *DefaultLinkService {
	return &DefaultLinkService{
		qrRepository: qrRepo,
		linksRepo:    linksRepo,
	}
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

	qr, err := s.qrRepository.GetQRCode(link, imgSize)
	if err != nil {
		return nil, err
	}
	return qr, nil
}
