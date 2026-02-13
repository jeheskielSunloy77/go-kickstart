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
		return false
	}
}
