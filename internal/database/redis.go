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
	godotenv.Load()

	redisURL := os.Getenv("UPSTASH_URL")
	if redisURL == "" {
		return fmt.Errorf("UPSTASH_URL not set in .env file")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Successfully connected to Upstash Redis!")
	return nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
		log.Println("Redis connection closed.")
	}
}
