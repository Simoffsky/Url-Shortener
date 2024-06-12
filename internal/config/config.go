package config

import "time"

type Config struct {
	Port        string
	HTTPTimeout time.Duration
	Name        string
}
