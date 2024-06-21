package stats

import (
	"url-shorter/internal/models"
	"url-shorter/internal/repository/stats"
)

type StatsService interface {
	GetStatForLink(short string) (*models.LinkStat, error)
	SendStat(stat *models.LinkStatVisitor) error
}

type StatsServiceDefault struct {
	statsRepo stats.StatsRepository
}

func NewStatsServiceDefault(statsRepo stats.StatsRepository) *StatsServiceDefault {
	return &StatsServiceDefault{
		statsRepo: statsRepo,
	}
}

func (s *StatsServiceDefault) GetStatForLink(short string) (*models.LinkStat, error) {
	return s.statsRepo.GetStatForLink(short)
}

func (s *StatsServiceDefault) SendStat(stat *models.LinkStatVisitor) error {
	panic("implement me (SendStat)")
}
