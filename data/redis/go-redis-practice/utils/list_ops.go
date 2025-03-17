package utils

import (
	"context"
	"go-redis-practice/config"
	"go-redis-practice/database"

	"github.com/redis/go-redis/v9"
)

type ListOps struct {
	Context     context.Context
	RedisClient *redis.Client
}

func NewListOps() *ListOps {
	redisConfig := config.NewRedisConfigFromEnv()
	redisClient := database.NewRedisClient(redisConfig)

	return &ListOps{
		RedisClient: redisClient.Client,
		Context:     context.Background(),
	}
}

// LPush inserts all the specified values at the head of the list stored at key.
func (l *ListOps) LPush(key string, values ...interface{}) (int64, error) {
	return l.RedisClient.LPush(l.Context, key, values...).Result()
}

// LRange returns the specified elements of the list stored at key.
func (l *ListOps) LRange(key string, start int64, stop int64) ([]string, error) {
	return l.RedisClient.LRange(l.Context, key, start, stop).Result()
}

// LLen returns the length of the list stored at key.
func (l *ListOps) LLen(key string) (int64, error) {
	return l.RedisClient.LLen(l.Context, key).Result()
}

// LPop removes and returns the first element of the list stored at key.
func (l *ListOps) LPop(key string) (string, error) {
	return l.RedisClient.LPop(l.Context, key).Result()
}

// RPush inserts all the specified values at the tail of the list stored at key.
func (l *ListOps) RPush(key string, values ...interface{}) (int64, error) {
	return l.RedisClient.RPush(l.Context, key, values...).Result()
}

// RPop removes and returns the last element of the list stored at key.
func (l *ListOps) RPop(key string) (string, error) {
	return l.RedisClient.RPop(l.Context, key).Result()
}

// LIndex returns the element at index index in the list stored at key.
func (l *ListOps) LIndex(key string, index int64) (string, error) {
	return l.RedisClient.LIndex(l.Context, key, index).Result()
}
