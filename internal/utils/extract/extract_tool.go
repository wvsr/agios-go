package utils

import (
	"agios/internal/prompts"
	"agios/internal/utils/llm"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type ToolType struct {
	ToolType string `json:"tool_type"`
}

func tryParseSearchToolOutput(raw string) (*ToolType, bool) {
	jsonStr, ok := extractJSONSegment(raw)

	if !ok {
		return nil, false
	}

	var parsed ToolType
	err := json.Unmarshal([]byte(jsonStr), &parsed)

	if err != nil {
		return nil, false
	}
	return &parsed, true
}

func ExtractToolType(ctx context.Context, text string) (*ToolType, error) {
	formattedPrompt, err := prompts.ToolDetectorPrompt.Format(map[string]any{"user_query": text})

	if err != nil {
		return nil, err
	}

	// Use the new GenerateFullResponse function with retry logic
	var result string
	var llmErr error // Use a different variable name for the LLM error
	for attempt := 0; attempt < 2; attempt++ {
		result, llmErr = llm.GenerateFullResponse(ctx, formattedPrompt, nil)
		if llmErr == nil {
			// Attempt to parse the result if LLM call was successful
			parsed, ok := tryParseSearchToolOutput(result)
			if ok {
				return parsed, nil
			}
			// If parsing fails, continue to the next attempt
		}
		// If LLM call failed, or parsing failed, continue to the next attempt
	}

	// If both attempts failed to produce a parsable result
	if llmErr != nil {
		return nil, fmt.Errorf("failed to generate text after multiple attempts: %w", llmErr)
	}
	return nil, errors.New("failed to parse LLM output into structured data after multiple attempts")
}
