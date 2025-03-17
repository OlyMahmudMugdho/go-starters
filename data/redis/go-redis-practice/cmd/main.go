package main

import (
	"context"
	"go-redis-practice/config"
	"go-redis-practice/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	redisConfig := config.NewRedisConfigFromEnv()
	redisClient := database.NewRedisClient(redisConfig)
	defer redisClient.Client.Close()

	ctx := context.Background()
	status := redisClient.Client.Ping(ctx)
	log.Println(status)
}
