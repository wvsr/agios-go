package llm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/strrl/tavily-go/pkg/tavily"
)

type SearchResponse struct {
	Results []tavily.SearchResult `json:"results"`
	Answer  string                `json:"answer"`
	Images  []string              `json:"images"`
}

func TavilySearch(query string) (SearchResponse, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load:", err)
	}

	tavilyApiKey := os.Getenv("TAVILY_API_KEY")
	if tavilyApiKey == "" {
		return SearchResponse{}, fmt.Errorf("TAVILY_API_KEY not set in environment")
	}

	client := tavily.NewClient(tavilyApiKey)

	resp, err := client.SearchWithOptions(
		context.Background(),
		query,
		tavily.WithMaxResults(10),
	)
	if err != nil {
		log.Printf("Tavily search failed: %v", err)
		return SearchResponse{}, err
	}

	formattedResponse := SearchResponse{
		Results: resp.Results,
		Answer:  resp.Answer,
		Images:  resp.Images,
	}

	return formattedResponse, nil
}
