package utils

func extractJSONSegment(raw string) (string, bool) {
	start := -1
	depth := 0

	for i, r := range raw {
		switch r {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && start != -1 {
				return raw[start : i+1], true
			}
		}
	}

	return "", false
}
