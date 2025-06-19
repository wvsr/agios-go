package llm

import (
	"context"
	"log"

	"github.com/strrl/tavily-go/pkg/tavily"
)

func TavilySearch(apiKey, query string) ([]tavily.SearchResult, error) {
	client := tavily.NewClient(apiKey)

	resp, err := client.SearchWithOptions(
		context.Background(),
		query,
		tavily.WithMaxResults(10),
	)

	if err != nil {
		log.Printf("tavily search failed: %v", err)
		return []tavily.SearchResult{}, err
	}

	return resp.Results, nil
}
