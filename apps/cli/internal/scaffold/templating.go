package scaffold

import "strings"

func ReplaceTokens(input string, replacements map[string]string) string {
	out := input
	for key, value := range replacements {
		out = strings.ReplaceAll(out, key, value)
	}
	return out
}
