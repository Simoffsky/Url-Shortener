package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func UnsetEnv() {
	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("QR_GRPC_ADDR")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("AUTH_GRPC_ADDR")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("DB_CONN")
	os.Unsetenv("STATS_ADDR")
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_GROUP")
	os.Unsetenv("KAFKA_TOPIC")

}
func TestNewEnvConfig(t *testing.T) {
	os.Setenv("PORT", "5000")
	os.Setenv("HTTP_TIMEOUT", "15")
	os.Setenv("LOGGER_LEVEL", "DEBUG")
	os.Setenv("QR_GRPC_ADDR", "localhost:50051")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("AUTH_GRPC_ADDR", "localhost:50052")
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("DB_CONN", "host=localhost port=5432 user=postgres password=postgres dbname=url_shorter sslmode=disable")
	os.Setenv("STATS_ADDR", "localhost:50053")
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("KAFKA_GROUP", "url_shorter")
	os.Setenv("KAFKA_TOPIC", "stats")

	expected := Config{
		ServerPort:   "5000",
		HTTPTimeout:  15 * time.Second,
		LoggerLevel:  "DEBUG",
		QrAddr:       "localhost:50051",
		RedisAddr:    "localhost:6379",
		AuthAddr:     "localhost:50052",
		JwtSecret:    "secret",
		DbConn:       "host=localhost port=5432 user=postgres password=postgres dbname=url_shorter sslmode=disable",
		StatsAddr:    "localhost:50053",
		KafkaBrokers: []string{"localhost:9092"},
		KafkaGroup:   "url_shorter",
		KafkaTopic:   "stats",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the expected values")

	UnsetEnv()
}

func TestNewEnvConfigWithDefaults(t *testing.T) {
	UnsetEnv()

	expected := Config{
		ServerPort:  "8080",
		HTTPTimeout: 10 * time.Second,
		LoggerLevel: "DEBUG",
		QrAddr:      "localhost:50051",
		RedisAddr:   "localhost:6379",
		AuthAddr:    "localhost:50052",
		JwtSecret:   "secret",
		DbConn:      "host=localhost port=5432 user=postgres password=postgres dbname=url_shorter sslmode=disable",
		StatsAddr:   "localhost:50053",
		KafkaBrokers: []string{
			"localhost:9092",
		},
		KafkaGroup: "url_shorter",
		KafkaTopic: "stats",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the default values")
}
