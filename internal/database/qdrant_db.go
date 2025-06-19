package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var QdrantClient qdrant.QdrantClient

func ConnectToQdrant() {
	godotenv.Load()

	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		log.Println("QDRANT_URL not set in .env file")
		return
	}

	qdrantAPIKey := os.Getenv("QDRANT_API_KEY")

	address := qdrantURL

	var dialOptions []grpc.DialOption

	if qdrantAPIKey != "" {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(apiKeyAuth{
			apiKey: qdrantAPIKey,
		}))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(address, dialOptions...)
	if err != nil {
		log.Println("QDRANT_URL not set in .env file")
		return
	}
	if err != nil {
		log.Printf("failed to connect to Qdrant: %v", err)
		return
	}

	QdrantClient = qdrant.NewQdrantClient(conn)

	log.Println("Successfully connected to Qdrant!")
}

func GetQdrantClient() qdrant.QdrantClient {
	return QdrantClient
}

type apiKeyAuth struct {
	apiKey string
}

func (a apiKeyAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"api-key": a.apiKey}, nil
}

func (a apiKeyAuth) RequireTransportSecurity() bool {
	return true
}
