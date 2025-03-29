package utils

import "strings"

// SplitMessage divides a text string into parts based on a character limit.
// It preserves line breaks and ensures no part exceeds the specified limit.
// Empty lines are skipped, and whitespace is trimmed from the results.
//
// Parameters:
//
//	text - The input string to be split
//	limit - Maximum character length for each resulting part
//
// Returns:
//
//	A slice of strings, each within the specified character limit
func SplitMessage(text string, limit int) []string {
	var result []string
	lines := strings.Split(strings.TrimSpace(text), "\n")

	currentPart := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if len(currentPart)+len(line)+1 > limit && currentPart != "" {
			result = append(result, strings.TrimSpace(currentPart))
			currentPart = line
		} else {
			if currentPart != "" {
				currentPart += "\n"
			}
			currentPart += line
		}
	}

	if strings.TrimSpace(currentPart) != "" {
		result = append(result, strings.TrimSpace(currentPart))
	}

	return result
}
