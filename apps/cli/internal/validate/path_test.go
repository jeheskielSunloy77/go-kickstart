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

func TestResolveProjectDestination_EmptyBaseUsesCwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	got, err := ResolveProjectDestination("", "demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(cwd, "demo")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveProjectDestination_RelativeBase(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	got, err := ResolveProjectDestination("some-dir", "demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(cwd, "some-dir", "demo")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveProjectDestination_AbsoluteBase(t *testing.T) {
	base := t.TempDir()
	got, err := ResolveProjectDestination(base, "demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(base, "demo")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveProjectDestination_SmartDedupe(t *testing.T) {
	base := filepath.Join(t.TempDir(), "demo")
	got, err := ResolveProjectDestination(base, "demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != base {
		t.Fatalf("expected %s, got %s", base, got)
	}
}

func TestResolveProjectDestination_SmartDedupeWithTrailingSlash(t *testing.T) {
	base := filepath.Join(t.TempDir(), "demo")
	got, err := ResolveProjectDestination(base+string(os.PathSeparator), "demo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != base {
		t.Fatalf("expected %s, got %s", base, got)
	}
}
