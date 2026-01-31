package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestScaffoldFromFS_StrictTemplates(t *testing.T) {
	tmp := t.TempDir()
	cfg := DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.ModulePath = "github.com/acme/demo"
	cfg.Destination = tmp
	cfg.InitGit = false

	fsys := fstest.MapFS{
		// Should NOT be templated (no .tmpl suffix)
		"plain.txt": &fstest.MapFile{Data: []byte("{{PROJECT_NAME}}")},
		// Should be templated + suffix stripped
		"templated.txt.tmpl": &fstest.MapFile{Data: []byte("{{PROJECT_NAME}}")},
	}

	if err := ScaffoldFromFS(cfg, true, fsys, nil); err != nil {
		t.Fatalf("scaffold failed: %v", err)
	}

	plain, err := os.ReadFile(filepath.Join(tmp, "plain.txt"))
	if err != nil {
		t.Fatalf("read plain.txt: %v", err)
	}
	if string(plain) != "{{PROJECT_NAME}}" {
		t.Fatalf("expected non-templated file to remain unchanged, got %q", string(plain))
	}

	templated, err := os.ReadFile(filepath.Join(tmp, "templated.txt"))
	if err != nil {
		t.Fatalf("read templated.txt: %v", err)
	}
	if string(templated) != "demo" {
		t.Fatalf("expected templated content, got %q", string(templated))
	}
}

