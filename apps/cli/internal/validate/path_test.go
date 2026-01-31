package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDestination_DefaultsToCwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	got, err := ResolveDestination("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != cwd {
		t.Fatalf("expected %s, got %s", cwd, got)
	}
}

func TestIsNonEmptyDir(t *testing.T) {
	dir := t.TempDir()
	nonEmpty, err := IsNonEmptyDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonEmpty {
		t.Fatalf("expected empty dir")
	}

	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	nonEmpty, err = IsNonEmptyDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !nonEmpty {
		t.Fatalf("expected non-empty dir")
	}
}
