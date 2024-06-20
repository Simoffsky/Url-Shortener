package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvConfig(t *testing.T) {
	os.Setenv("PORT", "5000")
	os.Setenv("HTTP_TIMEOUT", "15")
	os.Setenv("LOGGER_LEVEL", "DEBUG")
	os.Setenv("QR_ADDR", "localhost:50051")
	os.Setenv("REDIS_ADDR", "localhost:6379")

	expected := Config{
		ServerPort:  "5000",
		HTTPTimeout: 15 * time.Second,
		LoggerLevel: "DEBUG",
		QrAddr:      "localhost:50051",
		RedisAddr:   "localhost:6379",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the expected values")

	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("QR_ADDR")
	os.Unsetenv("REDIS_ADDR")
}

func TestNewEnvConfigWithDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("QR_ADDR")
	os.Unsetenv("REDIS_ADDR")

	expected := Config{
		ServerPort:  "8080",
		HTTPTimeout: 10 * time.Second,
		LoggerLevel: "DEBUG",
		QrAddr:      "localhost:50051",
		RedisAddr:   "localhost:6379",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the default values")
}
