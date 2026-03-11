package scaffold

import (
	"regexp"
	"strings"
)

var (
	includeWebBlockRe           = regexp.MustCompile(`(?s)\{\{IF_INCLUDE_WEB\}\}(.*?)\{\{END_IF_INCLUDE_WEB\}\}`)
	observabilityEnabledBlockRe = regexp.MustCompile(`(?s)\{\{IF_OBSERVABILITY_ENABLED\}\}(.*?)\{\{END_IF_OBSERVABILITY_ENABLED\}\}`)
	observabilityGrafanaBlockRe = regexp.MustCompile(`(?s)\{\{IF_OBSERVABILITY_GRAFANA\}\}(.*?)\{\{END_IF_OBSERVABILITY_GRAFANA\}\}`)
)

func ApplyTemplateConditions(input string, cfg ScaffoldConfiguration) string {
	if cfg.IncludeWeb {
		input = includeWebBlockRe.ReplaceAllString(input, "$1")
	} else {
		input = includeWebBlockRe.ReplaceAllString(input, "")
	}

	if cfg.Observability == ObservabilityGrafanaOSS {
		input = observabilityEnabledBlockRe.ReplaceAllString(input, "$1")
		input = observabilityGrafanaBlockRe.ReplaceAllString(input, "$1")
	} else {
		input = observabilityEnabledBlockRe.ReplaceAllString(input, "")
		input = observabilityGrafanaBlockRe.ReplaceAllString(input, "")
	}

	return input
}

func ReplaceTokens(input string, replacements map[string]string) string {
	out := input
	for key, value := range replacements {
		out = strings.ReplaceAll(out, key, value)
	}
	return out
}
