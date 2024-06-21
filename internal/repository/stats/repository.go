package stats

import "url-shorter/internal/models"

type StatsRepository interface {
	GetStatForLink(short string) (*models.LinkStat, error)
	AddStat(models.LinkStat) error
}
