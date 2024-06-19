package qr

import (
	"fmt"
	"net/url"
	"time"
	"url-shorter/internal/models"
	"url-shorter/internal/qr/cache"
	"url-shorter/pkg/log"

	"github.com/skip2/go-qrcode"
)

// image has imgSize x imgSize pixels
type QRService interface {
	GetQRCode(link string, size int) ([]byte, error)
}

type QRServiceDefault struct {
	logger log.Logger
	cache  cache.Cache
}

func NewDefaultQRService(cache cache.Cache, logger log.Logger) *QRServiceDefault {
	return &QRServiceDefault{
		cache:  cache,
		logger: logger,
	}
}

func (s *QRServiceDefault) GetQRCode(link string, imgSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return nil, models.ErrWrongLinkFormat
	}

	cacheKey := fmt.Sprintf("%s:%d", link, imgSize)

	cachedQR, err := s.cache.Get(cacheKey)
	if err == nil {
		s.logger.Debug("QR code found in cache")
		return cachedQR, nil
	}

	s.logger.Debug("QR code not found in cache, generating new one")
	rawPng, err := qrcode.Encode(link, qrcode.Medium, imgSize)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(cacheKey, rawPng, 100*time.Hour)
	if err != nil {
		s.logger.Error("Failed to save QR code to cache")
		return nil, err
	}

	return rawPng, nil
}
