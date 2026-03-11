package scaffold

import "strings"

func ShouldSkipForConfig(cfg ScaffoldConfiguration) func(path string) bool {
	return func(path string) bool {
		if !cfg.IncludeWeb && (path == "apps/web" || strings.HasPrefix(path, "apps/web/")) {
			return true
		}
		if !cfg.IncludeWeb && (path == "packages/ui" || strings.HasPrefix(path, "packages/ui/")) {
			return true
		}
		if !cfg.IncludeDocker && strings.HasPrefix(path, "docker-compose") {
			return true
		}
		if cfg.Observability != ObservabilityGrafanaOSS && (path == "ops/observability" || strings.HasPrefix(path, "ops/observability/")) {
			return true
		}
		if cfg.Observability != ObservabilityGrafanaOSS && (path == "apps/api/internal/infrastructure/observability" || strings.HasPrefix(path, "apps/api/internal/infrastructure/observability/")) {
			return true
		}
		return false
	}
}
