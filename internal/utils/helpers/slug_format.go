package helpers

import (
	"strings"
	"unicode"
)

func GetFirstNWords(text string, n int) string {
	words := []string{}
	inWord := false
	wordStart := 0

	for i, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				wordStart = i
				inWord = true
			}
		} else {
			if inWord {
				words = append(words, text[wordStart:i])
				if len(words) == n {
					break
				}
			}
			inWord = false
		}
	}
	if inWord && len(words) < n {
		words = append(words, text[wordStart:])
	}

	return strings.Join(words, " ")
}

func Slugify(text string) string {
	slug := strings.ToLower(text)
	slug = strings.ReplaceAll(slug, " ", "-")

	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = strings.Trim(result.String(), "-")

	return slug
}

func ValidateSlugFormat(slug, queryText string) bool {
	slugifiedFirstWords := Slugify(GetFirstNWords(queryText, 5))

	if queryText == "" {
		return false
	}

	if !strings.HasPrefix(slug, slugifiedFirstWords) {
		return false
	}

	nanoidPart := strings.TrimPrefix(slug, slugifiedFirstWords)
	nanoidPart = strings.TrimLeft(nanoidPart, "-")

	return nanoidPart != ""
}
