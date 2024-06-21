package models

type LinkStat struct {
	ShortLink      string `json:"short_link"`
	Clicks         int    `json:"clicks"`
	LastAccessedAt int64  `json:"last_accessed_at"` // Unix timestamp
	UniqueClicks   int    `json:"unique_clicks"`
}

type LinkStatVisitor struct {
	LinkShort string      `json:"link_short"`
	Visitor   LinkVisitor `json:"visitor"`
	TimeAt    int64       `json:"time_at"` // Unix timestamp
}
