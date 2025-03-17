package utils

import (
	"context"
	"go-redis-practice/config"
	"go-redis-practice/database"

	"github.com/redis/go-redis/v9"
)

type StringOps struct {
	RedisClient *redis.Client
	Context     context.Context
}

func NewStringOps() *StringOps {
	redisConfig := config.NewRedisConfigFromEnv()
	redisClient := database.NewRedisClient(redisConfig)

	return &StringOps{
		RedisClient: redisClient.Client,
		Context:     context.Background(),
	}
}

func (s *StringOps) Set(key string, value string) (string, error) {
	return s.RedisClient.Set(s.Context, key, value, 0).Result()
}

func (s *StringOps) Get(key string) (string, error) {
	return s.RedisClient.Get(s.Context, key).Result()
}
