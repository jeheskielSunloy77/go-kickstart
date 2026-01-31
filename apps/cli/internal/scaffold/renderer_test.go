package scaffold

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestRenderFS(t *testing.T) {
	fsys := fstest.MapFS{
		"README.md":           &fstest.MapFile{Data: []byte("hello")},
		"apps/web/index.html": &fstest.MapFile{Data: []byte("web")},
		"apps/api/main.go":    &fstest.MapFile{Data: []byte("api")},
		"docker-compose.yml":  &fstest.MapFile{Data: []byte("compose")},
		".git/HEAD":           &fstest.MapFile{Data: []byte("ref")},
	}

	tmp := t.TempDir()
	err := RenderFS(fsys, tmp, func(path string) bool {
		return path == ".git"
	}, nil)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "README.md")); err != nil {
		t.Fatalf("expected README.md to exist: %v", err)
	}
}
