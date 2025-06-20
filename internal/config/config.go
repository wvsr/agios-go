package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	DBSSLMode        string
	ExaAPIKey        string
	TavilyAPIKey     string
	CerebrasAPIKey   string
	BraveAPIKey      string
	Neo4jURI         string
	Neo4jUsername    string
	Neo4jPassword    string
	AuraInstanceID   string
	AuraInstanceName string
	GoogleAPIKey     string
	GoogleMapKey     string
	QdrantAPIKey     string
	QdrantURL        string
	UpstashAPIKey    string
	UpstashURL       string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBPort:           os.Getenv("DB_PORT"),
		DBSSLMode:        os.Getenv("DB_SSLMODE"),
		ExaAPIKey:        os.Getenv("EXA_API_KEY"),
		TavilyAPIKey:     os.Getenv("TAVILY_API_KEY"),
		CerebrasAPIKey:   os.Getenv("CEREBRAS_API_KEY"),
		BraveAPIKey:      os.Getenv("BRAVE_API_KEY"),
		Neo4jURI:         os.Getenv("NEO4J_URI"),
		Neo4jUsername:    os.Getenv("NEO4J_USERNAME"),
		Neo4jPassword:    os.Getenv("NEO4J_PASSWORD"),
		AuraInstanceID:   os.Getenv("AURA_INSTANCEID"),
		AuraInstanceName: os.Getenv("AURA_INSTANCENAME"),
		GoogleAPIKey:     os.Getenv("GOOGLE_API_KEY"),
		GoogleMapKey:     os.Getenv("GOOGLE_MAP_KEY"),
		QdrantAPIKey:     os.Getenv("QDRANT_API_KEY"),
		QdrantURL:        os.Getenv("QDRANT_URL"),
		UpstashAPIKey:    os.Getenv("UPSTASH_API_KEY"),
		UpstashURL:       os.Getenv("UPSTASH_URL"),
	}

	return cfg, nil
}
