package uttils

import (
	"unicode"
)

func WordCount(text string) int {
	inWord := false
	count := 0

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	return count
}
