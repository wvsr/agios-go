package database

import (
	"context"
	"fmt"
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
		return fmt.Errorf("DATABASE_URL not set in .env file")
	}

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return fmt.Errorf("failed to parse Neon database URL: %w", err)
	}

	NeonPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create Neon connection pool: %w", err)
	}

	err = NeonPool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to ping Neon database: %w", err)
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
