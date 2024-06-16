package repository

import (
	"sync"
	"url-shorter/internal/models"
)

type MemoryLinksRepository struct {
	mx    *sync.RWMutex
	links map[string]string
}

func NewMemoryLinksRepository() *MemoryLinksRepository {
	return &MemoryLinksRepository{
		mx:    &sync.RWMutex{},
		links: make(map[string]string),
	}
}

func (r *MemoryLinksRepository) CreateLink(url, short string) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, ok := r.links[short]
	if ok {
		return "", models.ErrLinkAlreadyExists
	}

	r.links[short] = url
	return short, nil
}

func (r *MemoryLinksRepository) GetLink(short string) (string, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	url, ok := r.links[short]
	if !ok {
		return "", models.ErrLinkNotFound
	}

	return url, nil
}
