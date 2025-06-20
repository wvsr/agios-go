package database

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
)

var QdrantClient *qdrant.Client

func ConnectToQdrant() error {
	_ = godotenv.Load()

	rawURL := os.Getenv("QDRANT_URL")
	apiKey := os.Getenv("QDRANT_API_KEY")

	if rawURL == "" {
		return fmt.Errorf("QDRANT_URL not set in environment")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid QDRANT_URL: %w", err)
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "6334" // default gRPC port for Qdrant Cloud
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}

	cfg := &qdrant.Config{
		Host:   host,
		Port:   portNum,
		APIKey: apiKey,
		UseTLS: true,
	}

	client, err := qdrant.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	QdrantClient = client
	log.Println("Successfully connected to Qdrant")

	return nil
}

func GetQdrantClient() *qdrant.Client {
	return QdrantClient
}
