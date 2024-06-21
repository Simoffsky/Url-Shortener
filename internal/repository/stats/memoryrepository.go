package stats

import (
	"errors"
	"sync"
	"url-shorter/internal/models"
)

// MemoryStatsRepository реализация интерфейса StatsRepository, хранящая данные в памяти
type MemoryStatsRepository struct {
	mu    *sync.RWMutex
	stats map[string]*models.LinkStat
}

func NewMemoryStatsRepository() *MemoryStatsRepository {
	return &MemoryStatsRepository{
		stats: make(map[string]*models.LinkStat),
	}
}

func (r *MemoryStatsRepository) GetStatForLink(short string) (*models.LinkStat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if stat, exists := r.stats[short]; exists {
		return stat, nil
	}
	return nil, errors.New("stat not found")
}

func (r *MemoryStatsRepository) AddStat(stat models.LinkStat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stats[stat.ShortLink] = &stat
	return nil
}
