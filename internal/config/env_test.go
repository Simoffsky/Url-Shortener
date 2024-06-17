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
	expected := Config{
		Port:        "5000",
		HTTPTimeout: 15 * time.Second,
		LoggerLevel: "DEBUG",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the expected values")

	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("LOGGER_LEVEL")
}

func TestNewEnvConfigWithDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("LOGGER_LEVEL")
	expected := Config{
		Port:        "8080",
		HTTPTimeout: 10 * time.Second,
		LoggerLevel: "DEBUG",
	}

	config := NewEnvConfig()

	assert.EqualValues(t, expected, config, "Config struct should match the default values")
}
