package repository

import (
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/sirupsen/logrus"
)

type RedisRepository struct {
	Redis *redis.Storage
	Log   *logrus.Logger
}

func NewRedisRepository(redis *redis.Storage, log *logrus.Logger) *RedisRepository {
	return &RedisRepository{
		Redis: redis,
		Log:   log,
	}
}

func (r *RedisRepository) Set(key string, value []byte, duration time.Duration) error {
	if r.Redis == nil {
		return nil
	}
	return r.Redis.Set(key, value, duration)
}

func (r *RedisRepository) Get(key string) ([]byte, error) {
	if r.Redis == nil {
		return nil, nil
	}
	return r.Redis.Get(key)
}

func (r *RedisRepository) Del(key string) error {
	if r.Redis == nil {
		return nil
	}
	return r.Redis.Delete(key)
}
