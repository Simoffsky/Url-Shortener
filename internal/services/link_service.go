package services

import (
	"url-shorter/internal/models"
	"url-shorter/internal/repository"
)

type LinkService interface {
	GetLink(short string) (*models.Link, error)
	RemoveLink(short string) error
	CreateLink(models.Link) error
}

type DefaultLinkService struct {
	linksRepo repository.LinksRepository
}

func NewDefaultLinkService(linksRepo repository.LinksRepository) *DefaultLinkService {
	return &DefaultLinkService{linksRepo: linksRepo}
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
