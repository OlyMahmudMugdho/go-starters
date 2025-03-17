package main

import (
	"go-redis-practice/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := server.NewServer()
	server.Start()
}
