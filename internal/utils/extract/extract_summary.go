package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"agios/internal/prompts"
	"agios/internal/utils/llm"
)

// KeyTakeaway represents a single takeaway.
type KeyTakeaway struct {
	Text            string  `json:"text"`
	ConfidenceScore float64 `json:"confidence_score"`
}

// Metric represents a single metric item.
type Metric struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

// ExtractionOutput is the full structured response.
type ExtractionOutput struct {
	KeyTakeaways       []KeyTakeaway `json:"key_takeaways"`
	RelatedSearchTerms []string      `json:"related_search_terms"`
	ShortSummary       string        `json:"short_summary"`
	Metrics            []Metric      `json:"metrics"`
}

// tryParseLLMOutput tries to parse the LLM output into ExtractionOutput.
func tryParseSummaryLLMOutput(raw string) (*ExtractionOutput, bool) {
	jsonStr, ok := extractJSONSegment(raw)
	if !ok {
		return nil, false
	}

	var parsed ExtractionOutput
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		return nil, false
	}

	return &parsed, true
}

// ExtractTakeaways sends input text to the LLM and tries to parse structured output.
func ExtractTakeaways(ctx context.Context, input string) (*ExtractionOutput, error) {
	prompt, err := prompts.SummaryPrompt.Format(map[string]any{
		"input_text": input,
	})
	if err != nil {
		return nil, err
	}

	// Use the new GenerateFullResponse function with retry logic
	var result string
	var llmErr error // Use a different variable name for the LLM error
	for attempt := 0; attempt < 2; attempt++ {
		result, llmErr = llm.GenerateFullResponse(ctx, prompt, nil)
		if llmErr == nil {
			// Attempt to parse the result if LLM call was successful
			parsed, ok := tryParseSummaryLLMOutput(result)
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
