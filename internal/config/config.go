package config

import (
	"fmt"
	"time"
)

type Config struct {
	ServerPort  string
	HTTPTimeout time.Duration
	LoggerLevel string

	RedisAddr string
	QrAddr    string

	AuthAddr  string
	JwtSecret string

	DbConn string

	StatsAddr string

	KafkaBrokers []string
	KafkaGroup   string
	KafkaTopic   string
}

func (c Config) String() string {

	return fmt.Sprintf("Port: %s, HTTPTimeout: %s, LoggerLevel: %s", c.ServerPort, c.HTTPTimeout, c.LoggerLevel)
}
