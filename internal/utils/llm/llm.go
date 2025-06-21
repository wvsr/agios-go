package llm

import (
	"context"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const DefualtModel = "gemini-1.5-flash"

func newClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create GenAI client: %w", err)
	}
	return client, nil
}

func createContentFromParts(query string, filePaths []string) ([]genai.Part, error) {
	var parts []genai.Part
	parts = append(parts, genai.Text(query))

	for _, filePath := range filePaths {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		mimeType := mime.TypeByExtension(filepath.Ext(filePath))
		if mimeType == "" {
			log.Printf("Warning: Could not determine MIME type for %s by extension, using application/octet-stream", filePath)
			mimeType = "application/octet-stream"
		}

		parts = append(parts, genai.ImageData(mimeType, data))
	}

	return parts, nil
}

func GenerateFullResponse(ctx context.Context, query string, filePaths []string) (string, error) {
	client, err := newClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	contents, err := createContentFromParts(query, filePaths)
	if err != nil {
		return "", fmt.Errorf("failed to create content parts: %w", err)
	}
	model := client.GenerativeModel(DefualtModel)

	resp, err := model.GenerateContent(ctx, contents...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		fullText := ""
		for _, part := range resp.Candidates[0].Content.Parts {
			fullText += fmt.Sprintf("%v", part)
		}
		return fullText, nil
	}

	return "", fmt.Errorf("no content generated")
}

func GenerateStreamResponse(ctx context.Context, query string, filePaths []string) (*genai.GenerateContentResponseIterator, error) {
	client, err := newClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	contents, err := createContentFromParts(query, filePaths)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create content parts: %w", err)
	}

	model := client.GenerativeModel(DefualtModel)

	iter := model.GenerateContentStream(ctx, contents...)

	// Return the iterator and the client. The caller is responsible for closing the client
	// after the iterator is fully consumed.
	// A better approach might be to wrap the iterator and client in a struct
	// that handles closing the client when the iterator is done.
	// For simplicity here, we return both and note the responsibility.
	// TODO: Consider wrapping iterator and client for better resource management.
	return iter, nil
}

// Example usage:
// iter, client, err := GenerateStreamResponse(...)
// defer client.Close() // Ensure client is closed after iterator is done
// for {
//     resp, err := iter.Next()
//     if err == iterator.Done {
//         break
//     }
//     if err != nil {
//         // Handle error
//     }
//     // Process resp
// }
