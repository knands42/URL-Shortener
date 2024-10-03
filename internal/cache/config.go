package cache

import (
	"knands42/url-shortener/internal/utils"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *utils.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
