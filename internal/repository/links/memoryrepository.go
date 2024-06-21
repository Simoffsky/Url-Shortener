package repository

import (
	"fmt"
	"sync"
	"url-shorter/internal/models"
)

type MemoryLinksRepository struct {
	mx    *sync.RWMutex
	links map[string]models.Link
}

func NewMemoryLinksRepository() *MemoryLinksRepository {
	return &MemoryLinksRepository{
		mx:    &sync.RWMutex{},
		links: make(map[string]models.Link),
	}
}

func (r *MemoryLinksRepository) CreateLink(link models.Link) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.links[link.ShortUrl]; ok {
		return fmt.Errorf("%w: %s", models.ErrLinkAlreadyExists, link.ShortUrl)
	}
	r.links[link.ShortUrl] = link

	return nil
}

func (r *MemoryLinksRepository) GetLink(short string) (*models.Link, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	link, ok := r.links[short]
	if !ok {
		return nil, fmt.Errorf("%w: %s", models.ErrLinkNotFound, short)
	}
	return &link, nil
}

func (r *MemoryLinksRepository) RemoveLink(short string) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.links[short]; !ok {
		return fmt.Errorf("%w: %s", models.ErrLinkNotFound, short)
	}
	delete(r.links, short)
	return nil
}

func (r *MemoryLinksRepository) EditLink(short string, editedLink models.Link) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.links[short]; !ok {
		return fmt.Errorf("%w: %s", models.ErrLinkNotFound, short)
	}
	r.links[short] = editedLink
	return nil
}
