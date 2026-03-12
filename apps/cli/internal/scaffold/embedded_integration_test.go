package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldProject_EmbeddedTemplateDefaultOutput(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "Demo App"
	cfg.ModulePath = "github.com/acme/demo-app"
	cfg.Destination = t.TempDir()
	cfg.InitGit = false

	if err := ScaffoldProject(cfg, false); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	assertExists(t, filepath.Join(cfg.Destination, "README.md"))
	assertExists(t, filepath.Join(cfg.Destination, "AGENTS.md"))
	assertExists(t, filepath.Join(cfg.Destination, "package.json"))
	assertExists(t, filepath.Join(cfg.Destination, "docker-compose.yml"))
	assertExists(t, filepath.Join(cfg.Destination, "apps/api/.env"))
	assertExists(t, filepath.Join(cfg.Destination, "apps/web/.env"))
	assertExists(t, filepath.Join(cfg.Destination, "apps/web"))
	assertExists(t, filepath.Join(cfg.Destination, "packages/ui"))

	readme := mustReadFile(t, filepath.Join(cfg.Destination, "README.md"))
	if !strings.Contains(readme, "demo-app/") {
		t.Fatalf("expected README to include kebab-cased project name")
	}

	apiEnv := mustReadFile(t, filepath.Join(cfg.Destination, "apps/api/.env"))
	if !strings.Contains(apiEnv, `API_PRIMARY.APP_NAME="Demo App"`) {
		t.Fatalf("expected api .env to include project name override")
	}
	if !strings.Contains(apiEnv, `API_DATABASE.NAME="app"`) {
		t.Fatalf("expected api .env to include default database name override")
	}

	rootPkg := mustReadFile(t, filepath.Join(cfg.Destination, "package.json"))
	if !strings.Contains(rootPkg, `"api:install": "cd apps/api && go mod tidy"`) {
		t.Fatalf("expected root package.json to install api dependencies with go mod tidy")
	}

	apiPkg := mustReadFile(t, filepath.Join(cfg.Destination, "apps/api/package.json"))
	if !strings.Contains(apiPkg, `"install": "go mod tidy"`) {
		t.Fatalf("expected api package.json to install dependencies with go mod tidy")
	}

	uiPkg := mustReadFile(t, filepath.Join(cfg.Destination, "packages/ui/package.json"))
	if !strings.Contains(uiPkg, `"tailwindcss": "^4.1.11"`) {
		t.Fatalf("expected ui package.json to declare tailwindcss")
	}
	if !strings.Contains(uiPkg, `"tw-animate-css": "^1.4.0"`) {
		t.Fatalf("expected ui package.json to declare tw-animate-css")
	}

	routerFile := mustReadFile(t, filepath.Join(cfg.Destination, "apps/web/src/router.tsx"))
	if !strings.Contains(routerFile, `from "react-router-dom"`) {
		t.Fatalf("expected web router to import from react-router-dom")
	}

	webEnv := mustReadFile(t, filepath.Join(cfg.Destination, "apps/web/.env"))
	if !strings.Contains(webEnv, `VITE_API_URL="http://localhost:8080"`) {
		t.Fatalf("expected web .env to be generated from example")
	}
}

func TestScaffoldProject_EmbeddedTemplateWithoutWebOrDocker(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = t.TempDir()
	cfg.IncludeWeb = false
	cfg.IncludeDocker = false
	cfg.InitGit = false

	if err := ScaffoldProject(cfg, false); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	assertNotExists(t, filepath.Join(cfg.Destination, "apps/web"))
	assertNotExists(t, filepath.Join(cfg.Destination, "packages/ui"))
	assertNotExists(t, filepath.Join(cfg.Destination, "docker-compose.yml"))

	readme := mustReadFile(t, filepath.Join(cfg.Destination, "README.md"))
	if strings.Contains(readme, "apps/web") {
		t.Fatalf("expected README to omit web references when web is disabled")
	}

	agents := mustReadFile(t, filepath.Join(cfg.Destination, "AGENTS.md"))
	if strings.Contains(agents, "apps/web") {
		t.Fatalf("expected AGENTS to omit web references when web is disabled")
	}

	pkg := mustReadFile(t, filepath.Join(cfg.Destination, "package.json"))
	if strings.Contains(pkg, "web:test") {
		t.Fatalf("expected package.json to omit web test script")
	}
	if strings.Contains(pkg, "ui:shadcn:add") {
		t.Fatalf("expected package.json to omit ui script")
	}
}

func TestScaffoldProject_EmbeddedTemplateWithGrafanaObservability(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = t.TempDir()
	cfg.InitGit = false
	cfg.Observability = ObservabilityGrafanaOSS

	if err := ScaffoldProject(cfg, false); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	assertExists(t, filepath.Join(cfg.Destination, "ops/observability/grafana/otel-collector.yml"))
	assertExists(t, filepath.Join(cfg.Destination, "apps/api/internal/infrastructure/observability"))

	readme := mustReadFile(t, filepath.Join(cfg.Destination, "README.md"))
	if !strings.Contains(readme, "docker compose --profile observability up --build") {
		t.Fatalf("expected README to include observability instructions")
	}

	apiEnv := mustReadFile(t, filepath.Join(cfg.Destination, "apps/api/.env"))
	if !strings.Contains(apiEnv, `OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"`) {
		t.Fatalf("expected api .env to include OTEL configuration")
	}
}

func TestScaffoldProject_EmbeddedTemplateInitGit(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = t.TempDir()
	cfg.IncludeWeb = false
	cfg.IncludeDocker = false

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "Go Kickstart Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Go Kickstart Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
	t.Setenv("GIT_TERMINAL_PROMPT", "0")

	if err := ScaffoldProject(cfg, false); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	assertExists(t, filepath.Join(cfg.Destination, ".git"))

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = cfg.Destination
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected git commit to exist: %v\n%s", err, string(output))
	}
	if strings.TrimSpace(string(output)) == "" {
		t.Fatalf("expected git rev-parse HEAD to return a commit hash")
	}
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(data)
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s not to exist", path)
	}
}
