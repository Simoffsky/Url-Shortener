package services

import "url-shorter/internal/models"

type StatService interface {
	GetStatForLink(short string) (*models.LinkStat, error)
	SendStat(stat *models.LinkStatVisitor) error
}

type StatServiceGRPC struct {
	statsClient
}
