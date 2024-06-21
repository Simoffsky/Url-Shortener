package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func NewEnvConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING:", err) // not returning error because we can work without .env file
	}
	config := Config{
		ServerPort:  getEnv("PORT", "8080"),
		HTTPTimeout: time.Duration(getEnvAsInt("HTTP_TIMEOUT", 10)) * time.Second,
		LoggerLevel: getEnv("LOGGER_LEVEL", "DEBUG"),
		QrAddr:      getEnv("QR_ADDR", "localhost:50051"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		AuthAddr:    getEnv("AUTH_GRPC_ADDR", "localhost:50052"),
		JwtSecret:   getEnv("JWT_SECRET", "secret"),
		DbConn:      getEnv("DB_CONN", "host=localhost port=5432 user=postgres password=postgres dbname=url_shorter sslmode=disable"),

		StatsAddr:    getEnv("STATS_ADDR", "localhost:50053"),
		KafkaBrokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
		KafkaGroup:   getEnv("KAFKA_GROUP", "url_shorter"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "stats"),
	}

	return config
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
