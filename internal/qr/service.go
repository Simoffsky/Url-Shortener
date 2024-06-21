package qr

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	qrChan := make(chan []byte, 1)
	errChan := make(chan error, 1)
	go func() {
		rawPng, err := qrcode.Encode(link, qrcode.Medium, imgSize)
		if err != nil {
			errChan <- err
			return
		}
		qrChan <- rawPng
	}()

	select {
	case rawPng := <-qrChan:
		err = s.cache.Set(cacheKey, rawPng, 100*time.Hour)
		if err != nil {
			s.logger.Error("Failed to save QR code to cache")
			return nil, err
		}
		return rawPng, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
