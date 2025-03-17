package utils

import (
	"context"
	"fmt"
	"go-redis-practice/config"
	"go-redis-practice/database"

	"github.com/redis/go-redis/v9"
)

type HashOps struct {
	RedisClient *redis.Client
	Context     context.Context
}

func NewHashOps() *HashOps {
	redisConfig := config.NewRedisConfigFromEnv()
	redisClient := database.NewRedisClient(redisConfig)

	return &HashOps{
		RedisClient: redisClient.Client,
		Context:     context.Background(),
	}
}

// HSet sets field in the hash stored at key to value.
func (h *HashOps) HSet(key string, values ...interface{}) (int64, error) {
	res, err := h.RedisClient.HSet(h.Context, key, values...).Result()
	if err != nil {
		fmt.Println("error", err)
	}

	return res, err
}

// HGet returns the value associated with field in the hash stored at key.
func (h *HashOps) HGet(key string, field string) (string, error) {
	return h.RedisClient.HGet(h.Context, key, field).Result()
}

// / HDel deletes a hash field.
func (h *HashOps) HDel(key string, fields ...string) (int64, error) {
	return h.RedisClient.HDel(h.Context, key, fields...).Result()
}

// HGetAll returns all fields and values of the hash stored at key.
func (h *HashOps) HGetAll(key string) (map[string]string, error) {
	return h.RedisClient.HGetAll(h.Context, key).Result()
}

// HExists returns whether a hash field exists or not.
func (h *HashOps) HExists(key string, field string) (bool, error) {
	return h.RedisClient.HExists(h.Context, key, field).Result()
}

// HKeys returns all field names in the hash stored at key.
func (h *HashOps) HKeys(key string) ([]string, error) {
	return h.RedisClient.HKeys(h.Context, key).Result()
}

// HVals returns all values in the hash stored at key.
func (h *HashOps) HVals(key string) ([]string, error) {
	return h.RedisClient.HVals(h.Context, key).Result()
}

// HLen returns the number of fields in the hash stored at key.
func (h *HashOps) HLen(key string) (int64, error) {
	return h.RedisClient.HLen(h.Context, key).Result()
}

// HIncrBy increments the integer value of a hash field by the given number.
func (h *HashOps) HIncrBy(key string, field string, increment int64) (int64, error) {
	return h.RedisClient.HIncrBy(h.Context, key, field, increment).Result()
}

// HIncrByFloat increments the float value of a hash field by the given amount.
func (h *HashOps) HIncrByFloat(key string, field string, increment float64) (float64, error) {
	return h.RedisClient.HIncrByFloat(h.Context, key, field, increment).Result()
}

// / HScan scans the hash stored at key for fields matching the given pattern.
func (h *HashOps) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return h.RedisClient.HScan(h.Context, key, cursor, match, count).Result()
}

// HSetNX sets the field in the hash stored at key to value, only if field does not yet exist.
func (h *HashOps) HSetNX(key string, field string, value string) (bool, error) {
	return h.RedisClient.HSetNX(h.Context, key, field, value).Result()
}

// HMSet is a convenience method to set multiple fields in a hash.
func (h *HashOps) HMSet(key string, fields map[string]interface{}) (bool, error) {
	return h.RedisClient.HMSet(h.Context, key, fields).Result()
}
