package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis() error {
	_ = godotenv.Load()

	redisURL := os.Getenv("UPSTASH_URL")
	if redisURL == "" {
		return fmt.Errorf("UPSTASH_URL not set in environment")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	RedisClient = redis.NewClient(opt)

	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Successfully connected to Redis!")

	return nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			fmt.Println("failed to close Redis:", err)
		}
		fmt.Println("Redis connection closed.")
	}
}
