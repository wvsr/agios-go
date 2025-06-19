package llm

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func GenerateText(ctx context.Context, prompt string) (string, string) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	defaultModel := "gemini-2.5-flash-preview-05-20"
	if apiKey == "" {
		log.Println("Error: GOOGLE_API_KEY not set")
		return "", ""
	}

	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(defaultModel))
	if err != nil {
		log.Printf("Error creating LLM: %v", err)
		return "", ""
	}

	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Printf("Error generating text: %v", err)
		return "", ""
	}

	return answer, ""
}
