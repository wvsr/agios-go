package utils

import (
	"agios/internal/prompts"
	"agios/internal/utils/llm"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type SearchTerm struct {
	SearchTerm []string `json:"search_term"`
}

func tryParseSearchTermsLLMOutput(raw string) (*SearchTerm, bool) {
	jsonStr, ok := extractJSONSegment(raw)
	if !ok {
		return nil, false
	}

	var parsed SearchTerm
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, false
	}

	return &parsed, true
}

func ExtractSearchTerms(ctx context.Context, text string) (*SearchTerm, error) {
	formattedPrompt, err := prompts.Search_term_prompt.Format(map[string]any{"inputText": text})
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
			parsed, ok := tryParseSearchTermsLLMOutput(result)
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
