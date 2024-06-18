package qr

import (
	"net/url"
	"sync"
	"url-shorter/internal/models"

	"github.com/skip2/go-qrcode"
)

// image has imgSize x imgSize pixels
type QRService interface {
	CreateQRCode(url string, imgSize int) ([]byte, error)
	GetQRCode(url string) ([]byte, error)
	DeleteQRCode(url string) error
}

type QRServiceDefault struct {
	mx        *sync.RWMutex
	QRStorage map[string][]byte
}

func NewDefaultQRService() *QRServiceDefault {
	return &QRServiceDefault{
		mx:        &sync.RWMutex{},
		QRStorage: map[string][]byte{},
	}
}

// returns raw png bytes of the QR code
// TODO: move saving to repository
func (s *QRServiceDefault) CreateQRCode(link string, imgSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return nil, models.ErrWrongLinkFormat
	}
	rawPng, err := qrcode.Encode(link, qrcode.Medium, imgSize)

	if err != nil {
		return nil, err
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	s.QRStorage[link] = rawPng
	return rawPng, nil
}

func (s *QRServiceDefault) GetQRCode(link string) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if qr, ok := s.QRStorage[link]; ok {
		return qr, nil
	}

	return nil, models.ErrLinkNotFound
}

func (s *QRServiceDefault) DeleteQRCode(link string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if _, ok := s.QRStorage[link]; ok {
		delete(s.QRStorage, link)
		return nil
	}

	return models.ErrLinkNotFound
}
