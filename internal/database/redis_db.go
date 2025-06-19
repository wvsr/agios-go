package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectToRedis() {
	godotenv.Load()

	redisURL := os.Getenv("UPSTASH_URL")
	if redisURL == "" {
		log.Println("UPSTASH_URL not set in .env file")
		return
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Println("failed to parse Redis URL:", err)
		return
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("failed to ping Redis:", err)
		return
	}

	log.Println("Successfully connected to Upstash Redis!")
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Println("failed to close Redis:", err)
		}
		log.Println("Redis connection closed.")
	}
}
