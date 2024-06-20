package repository

import "url-shorter/internal/models"

type LinksRepository interface {
	CreateLink(models.Link) error
	GetLink(short string) (*models.Link, error)
	EditLink(short string, editedLink models.Link) error
	RemoveLink(short string) error
}

type QrRepository interface {
	GetQRCode(link string, imgSize int) ([]byte, error)
}
