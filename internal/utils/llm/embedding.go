package llm

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms/googleai"
)

// CreateEmbedding creates embeddings for the given texts using the Gemini API.
func CreateEmbedding(ctx context.Context, texts []string) ([][]float32, string) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Println("GOOGLE_API_KEY not set")
		return nil, ""
	}

	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("error creating LLM: %v", err)
		return nil, ""
	}

	emb, err := llm.CreateEmbedding(ctx, texts)
	if err != nil {
		log.Printf("error creating embeddings: %v", err)
		return nil, ""
	}

	return emb, ""
}
