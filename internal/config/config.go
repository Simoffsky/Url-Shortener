package config

import (
	"fmt"
	"time"
)

type Config struct {
	Port        string
	HTTPTimeout time.Duration
	LoggerLevel string
}

func (c Config) String() string {
	return fmt.Sprintf("Port: %s, HTTPTimeout: %s, LoggerLevel: %s", c.Port, c.HTTPTimeout, c.LoggerLevel)
}
