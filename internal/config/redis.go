package config

import (
	"github.com/gofiber/storage/redis/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper, log *logrus.Logger) (storage *redis.Storage) {
	host := viper.GetString("database.redis.host")
	port := viper.GetInt("database.redis.port")
	password := viper.GetString("database.redis.password")

	log.Info("Connecting to redis")

	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Redis connection failed (non-critical): %+v", r)
			storage = nil
		}
	}()

	storage = redis.New(redis.Config{
		Host:     host,
		Port:     port,
		Password: password,
	})
	return
}
