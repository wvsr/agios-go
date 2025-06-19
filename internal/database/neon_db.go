package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var NeonPool *pgxpool.Pool

func ConnectToNeonDB() error {
	godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Println("DATABASE_URL not set in .env file")
		return nil
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Printf("failed to parse Neon database URL: %v", err)
		return nil
	}

	NeonPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("failed to create Neon connection pool: %v", err)
		return nil
	}

	err = NeonPool.Ping(context.Background())
	if err != nil {
		log.Printf("failed to ping Neon database: %v", err)
		return nil
	}

	log.Println("Successfully connected to NeonDB!")
	return nil
}

func GetNeonPool() *pgxpool.Pool {
	return NeonPool
}

func CloseNeonDB() {
	if NeonPool != nil {
		NeonPool.Close()
		log.Println("NeonDB connection closed.")
	}
}
