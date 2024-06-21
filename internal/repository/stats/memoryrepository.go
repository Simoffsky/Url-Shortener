package stats

import (
	"sync"
	"url-shorter/internal/models"
)

type MemoryStatsRepository struct {
	mu       *sync.RWMutex
	stats    map[string]*models.LinkStat
	visitors map[string]map[models.LinkVisitor]struct{} // short link -> visitors
}

func NewMemoryStatsRepository() *MemoryStatsRepository {
	return &MemoryStatsRepository{
		mu:       &sync.RWMutex{},
		stats:    make(map[string]*models.LinkStat),
		visitors: make(map[string]map[models.LinkVisitor]struct{}),
	}
}

func (r *MemoryStatsRepository) GetStatForLink(short string) (*models.LinkStat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if stat, exists := r.stats[short]; exists {
		return stat, nil
	}
	stat := &models.LinkStat{
		ShortLink:      short,
		Clicks:         0,
		UniqueClicks:   0,
		LastAccessedAt: 0,
	}
	r.stats[short] = stat
	return stat, nil
}

func (r *MemoryStatsRepository) AddStat(linkVisitor models.LinkStatVisitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.stats[linkVisitor.LinkShort]
	if !exists {
		stat := &models.LinkStat{
			ShortLink:      linkVisitor.LinkShort,
			Clicks:         0,
			UniqueClicks:   0,
			LastAccessedAt: 0,
		}
		r.stats[linkVisitor.LinkShort] = stat
	}

	visitors, exists := r.visitors[linkVisitor.LinkShort]
	if !exists {
		visitors = make(map[models.LinkVisitor]struct{})
		r.visitors[linkVisitor.LinkShort] = visitors
	}

	_, visitorExists := visitors[linkVisitor.Visitor]
	if !visitorExists {
		r.stats[linkVisitor.LinkShort].UniqueClicks++
		visitors[linkVisitor.Visitor] = struct{}{}
	}

	r.stats[linkVisitor.LinkShort].Clicks++
	r.stats[linkVisitor.LinkShort].LastAccessedAt = linkVisitor.TimeAt

	return nil
}
