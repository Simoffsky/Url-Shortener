package repository

import "url-shorter/internal/models"

type LinksRepository interface {
	CreateLink(models.Link) error
	GetLink(short string) (*models.Link, error)
	RemoveLink(short string) error
}
