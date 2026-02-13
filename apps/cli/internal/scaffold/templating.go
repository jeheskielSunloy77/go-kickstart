package scaffold

import (
	"regexp"
	"strings"
)

var includeWebBlockRe = regexp.MustCompile(`(?s)\{\{IF_INCLUDE_WEB\}\}(.*?)\{\{END_IF_INCLUDE_WEB\}\}`)

func ApplyTemplateConditions(input string, cfg ScaffoldConfiguration) string {
	if cfg.IncludeWeb {
		return includeWebBlockRe.ReplaceAllString(input, "$1")
	}
	return includeWebBlockRe.ReplaceAllString(input, "")
}

func ReplaceTokens(input string, replacements map[string]string) string {
	out := input
	for key, value := range replacements {
		out = strings.ReplaceAll(out, key, value)
	}
	return out
}
