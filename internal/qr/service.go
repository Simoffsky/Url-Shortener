package qr

import (
	"net/url"
	"sync"
	"url-shorter/internal/models"

	repository "url-shorter/internal/repository/qr"
)

// image has imgSize x imgSize pixels
type QRService interface {
	CreateQRCode(url string, imgSize int) ([]byte, error)
	GetQRCode(url string) ([]byte, error)
	DeleteQRCode(url string) error
}

type QRServiceDefault struct {
	mx   *sync.RWMutex
	repo repository.QrRepository
}

func NewDefaultQRService() *QRServiceDefault {
	return &QRServiceDefault{
		mx:   &sync.RWMutex{},
		repo: repository.NewMemoryQrRepository(),
	}
}

func (s *QRServiceDefault) CreateQRCode(link string, imgSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return nil, models.ErrWrongLinkFormat
	}
	qr, err := s.repo.CreateQRCode(link, imgSize)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (s *QRServiceDefault) GetQRCode(link string) ([]byte, error) {
	return s.repo.GetQRCode(link)
}

func (s *QRServiceDefault) DeleteQRCode(link string) error {
	return s.repo.DeleteQRCode(link)
}
