package utils

import (
	"agios/internal/prompts"
	"agios/internal/utils/llm"
	"context"
	"encoding/json"
	"errors"
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

	for range 2 {
		result, errMsg := llm.GenerateText(ctx, formattedPrompt)
		if errMsg != "" {
			continue
		}

		parsed, ok := tryParseSearchToolOutput(result)

		if ok {
			return parsed, nil
		}
	}

	return nil, errors.New("failed to parse LLM output into")
}
