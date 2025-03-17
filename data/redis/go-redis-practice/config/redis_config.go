package config

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisConfigFromEnv() *RedisConfig {
	host := os.Getenv("REDIS_HOST")
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return &RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}
