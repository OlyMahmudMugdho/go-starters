package utils

import (
	"context"
	"go-redis-practice/config"
	"go-redis-practice/database"

	"github.com/redis/go-redis/v9"
)

type PingUtils struct {
	RedisClient *redis.Client
}

func NewPingUtils() *PingUtils {
	redisConfig := config.NewRedisConfigFromEnv()
	redisClient := database.NewRedisClient(redisConfig)

	return &PingUtils{
		RedisClient: redisClient.Client,
	}
}

func (p *PingUtils) Ping() *redis.StatusCmd {
	defer p.RedisClient.Close()
	return p.RedisClient.Ping(context.Background())
}
