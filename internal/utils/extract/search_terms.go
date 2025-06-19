package utils

import (
	"agios/internal/prompts"
	"agios/internal/utils/llm"
	"context"
	"encoding/json"
	"errors"
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
	formattedPrompt, err := prompts.Search_term_prompt.Format(map[string]any{text: text})
	if err != nil {
		return nil, err
	}

	for range 2 {
		result, errMsg := llm.GenerateText(ctx, formattedPrompt)
		if errMsg != "" {
			continue
		}
		parsed, ok := tryParseSearchTermsLLMOutput(result)
		if ok {
			return parsed, nil
		}
	}

	return nil, errors.New("failed to parse LLM output into structured data")
}
