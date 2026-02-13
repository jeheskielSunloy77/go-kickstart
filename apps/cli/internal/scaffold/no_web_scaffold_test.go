package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestScaffoldFromFS_NoWebOutput(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = t.TempDir()
	cfg.IncludeWeb = false
	cfg.IncludeDocker = true
	cfg.InitGit = false

	fsys := fstest.MapFS{
		"README.md.tmpl":           &fstest.MapFile{Data: []byte("intro\n{{IF_INCLUDE_WEB}}apps/web\n## Web (apps/web)\n{{END_IF_INCLUDE_WEB}}")},
		"AGENTS.md.tmpl":           &fstest.MapFile{Data: []byte("ctx\n{{IF_INCLUDE_WEB}}App #2: Web (apps/web)\n{{END_IF_INCLUDE_WEB}}")},
		"package.json.tmpl":        &fstest.MapFile{Data: []byte("{\n\"scripts\": {\n\"api:test\": \"x\",\n{{IF_INCLUDE_WEB}}\"ui:shadcn:add\": \"z\",\n\"web:test\": \"y\",\n{{END_IF_INCLUDE_WEB}}\"ui\": \"z\"\n}\n}")},
		"docker-compose.yml.tmpl":  &fstest.MapFile{Data: []byte("services:\n  api: {}\n{{IF_INCLUDE_WEB}}  web: {}\n{{END_IF_INCLUDE_WEB}}")},
		"apps/web/index.html.tmpl": &fstest.MapFile{Data: []byte("web")},
		"packages/ui/package.json": &fstest.MapFile{Data: []byte("ui")},
		"apps/api/main.go":         &fstest.MapFile{Data: []byte("api")},
	}

	if err := ScaffoldFromFS(cfg, true, fsys, nil); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(cfg.Destination, "apps/web")); !os.IsNotExist(err) {
		t.Fatalf("expected apps/web directory not to exist")
	}
	if _, err := os.Stat(filepath.Join(cfg.Destination, "packages/ui")); !os.IsNotExist(err) {
		t.Fatalf("expected packages/ui directory not to exist")
	}

	readme, err := os.ReadFile(filepath.Join(cfg.Destination, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if strings.Contains(string(readme), "apps/web") {
		t.Fatalf("expected README.md without apps/web references")
	}

	agents, err := os.ReadFile(filepath.Join(cfg.Destination, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if strings.Contains(string(agents), "apps/web") {
		t.Fatalf("expected AGENTS.md without apps/web references")
	}

	pkg, err := os.ReadFile(filepath.Join(cfg.Destination, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	if strings.Contains(string(pkg), "web:") {
		t.Fatalf("expected package.json without web scripts")
	}
	if strings.Contains(string(pkg), "ui:shadcn:add") {
		t.Fatalf("expected package.json without ui script when web is excluded")
	}

	compose, err := os.ReadFile(filepath.Join(cfg.Destination, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("read docker-compose.yml: %v", err)
	}
	if strings.Contains(string(compose), "\n  web:") {
		t.Fatalf("expected docker-compose.yml without web service")
	}
}

func TestScaffoldFromFS_WithWebOutput(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = t.TempDir()
	cfg.IncludeWeb = true
	cfg.IncludeDocker = true
	cfg.InitGit = false

	fsys := fstest.MapFS{
		"README.md.tmpl":           &fstest.MapFile{Data: []byte("{{IF_INCLUDE_WEB}}apps/web{{END_IF_INCLUDE_WEB}}")},
		"AGENTS.md.tmpl":           &fstest.MapFile{Data: []byte("{{IF_INCLUDE_WEB}}App #2: Web (apps/web){{END_IF_INCLUDE_WEB}}")},
		"package.json.tmpl":        &fstest.MapFile{Data: []byte("{{IF_INCLUDE_WEB}}\"ui:shadcn:add\": \"z\",\"web:test\": \"y\"{{END_IF_INCLUDE_WEB}}")},
		"docker-compose.yml.tmpl":  &fstest.MapFile{Data: []byte("services:\n{{IF_INCLUDE_WEB}}  web: {}\n{{END_IF_INCLUDE_WEB}}")},
		"apps/web/index.html.tmpl": &fstest.MapFile{Data: []byte("web")},
		"packages/ui/package.json": &fstest.MapFile{Data: []byte("ui")},
	}

	if err := ScaffoldFromFS(cfg, true, fsys, nil); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(cfg.Destination, "apps/web")); err != nil {
		t.Fatalf("expected apps/web directory to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cfg.Destination, "packages/ui")); err != nil {
		t.Fatalf("expected packages/ui directory to exist: %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(cfg.Destination, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !strings.Contains(string(readme), "apps/web") {
		t.Fatalf("expected README.md to include apps/web references")
	}

	agents, err := os.ReadFile(filepath.Join(cfg.Destination, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(agents), "apps/web") {
		t.Fatalf("expected AGENTS.md to include apps/web references")
	}

	pkg, err := os.ReadFile(filepath.Join(cfg.Destination, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	if !strings.Contains(string(pkg), "web:test") {
		t.Fatalf("expected package.json to include web scripts")
	}
	if !strings.Contains(string(pkg), "ui:shadcn:add") {
		t.Fatalf("expected package.json to include ui script")
	}

	compose, err := os.ReadFile(filepath.Join(cfg.Destination, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("read docker-compose.yml: %v", err)
	}
	if !strings.Contains(string(compose), "\n  web:") {
		t.Fatalf("expected docker-compose.yml to include web service")
	}
}
