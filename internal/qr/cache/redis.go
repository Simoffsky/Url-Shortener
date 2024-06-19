package cache

import (
	"time"
	"url-shorter/internal/models"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	redisClient *redis.Client
}

func NewRedisCache(redisClient *redis.Client) *RedisCache {
	return &RedisCache{
		redisClient: redisClient,
	}
}

func (c *RedisCache) Get(key string) ([]byte, error) {
	bytes, err := c.redisClient.Get(key).Bytes()
	if err == redis.Nil {
		return nil, models.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return c.redisClient.Set(key, value, ttl).Err()
}
