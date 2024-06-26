package services

import (
	"fmt"
	"time"
	"url-shorter/internal/models"
	repository "url-shorter/internal/repository/links"
)

type LinkService interface {
	GetLink(short string) (*models.Link, error)
	RemoveLink(creatorLogin string, short string) error
	EditLink(creatorLogin string, short string, editedLink models.Link) error
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

	link, err := s.linksRepo.GetLink(short)
	if err != nil {
		return nil, err
	}
	fmt.Println(link)
	if link.ExpiredAt != 0 && link.ExpiredAt < time.Now().Unix() {
		return nil, models.ErrLinkExpired
	}
	return link, nil
}

func (s *DefaultLinkService) RemoveLink(creatorLogin string, short string) error {

	link, err := s.linksRepo.GetLink(short)
	if err != nil {
		return err
	}

	if link.CreatorLogin == "" || link.CreatorLogin != creatorLogin {
		return models.ErrForbidden
	}
	return s.linksRepo.RemoveLink(short)
}

func (s *DefaultLinkService) CreateLink(link models.Link) error {
	if err := link.Validate(); err != nil {
		return err
	}

	return s.linksRepo.CreateLink(link)
}

func (s *DefaultLinkService) EditLink(creatorLogin string, short string, editedLink models.Link) error {
	if err := editedLink.Validate(); err != nil {
		return err
	}
	link, err := s.GetLink(short)
	if err != nil {
		return err
	}

	if link.CreatorLogin == "" || link.CreatorLogin != creatorLogin {
		return models.ErrForbidden
	}

	return s.linksRepo.EditLink(short, editedLink)
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
