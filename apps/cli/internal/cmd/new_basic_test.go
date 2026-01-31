package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func TestScaffoldBasicFlow(t *testing.T) {
	tmp := t.TempDir()
	cfg := scaffold.DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = tmp
	cfg.IncludeWeb = true
	cfg.IncludeDocker = true
	cfg.InitGit = false

	fsys := fstest.MapFS{
		"README.md.tmpl": &fstest.MapFile{Data: []byte("{{PROJECT_NAME}}")},
	}

	if err := scaffold.ScaffoldFromFS(cfg, true, fsys, nil); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "README.md"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(data) != "demo" {
		t.Fatalf("expected templated content, got %s", string(data))
	}
}
