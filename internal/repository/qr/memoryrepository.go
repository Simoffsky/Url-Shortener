package qr

import (
	"sync"
	"url-shorter/internal/models"

	"github.com/skip2/go-qrcode"
)

type MemoryQrRepository struct {
	mx      *sync.RWMutex
	qrCodes map[string][]byte
}

func NewMemoryQrRepository() *MemoryQrRepository {
	return &MemoryQrRepository{
		mx:      &sync.RWMutex{},
		qrCodes: make(map[string][]byte),
	}
}

func (r *MemoryQrRepository) CreateQRCode(link string, imgSize int) ([]byte, error) {
	rawPng, err := qrcode.Encode(link, qrcode.Medium, imgSize)

	if err != nil {
		return nil, err
	}

	r.mx.Lock()
	defer r.mx.Unlock()
	r.qrCodes[link] = rawPng
	return rawPng, nil
}

func (r *MemoryQrRepository) GetQRCode(link string) ([]byte, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	if qr, ok := r.qrCodes[link]; ok {
		return qr, nil
	}

	return nil, models.ErrLinkNotFound
}

func (r *MemoryQrRepository) DeleteQRCode(link string) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.qrCodes[link]; ok {
		delete(r.qrCodes, link)
		return nil
	}

	return models.ErrLinkNotFound
}
