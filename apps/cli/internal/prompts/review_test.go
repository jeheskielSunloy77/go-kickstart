package prompts

import (
	"path/filepath"
	"testing"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func TestResolveDisplayDestinationUsesResolvedProjectPath(t *testing.T) {
	cfg := scaffold.DefaultConfig()
	cfg.ProjectName = "my-app"
	cfg.Destination = ""

	got := resolveDisplayDestination(cfg)
	want, err := filepath.Abs(filepath.Join(".", "my-app"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveDisplayDestinationResolvesProvidedBasePath(t *testing.T) {
	cfg := scaffold.DefaultConfig()
	cfg.ProjectName = "demo"
	cfg.Destination = "somewhere"

	got := resolveDisplayDestination(cfg)
	want, err := filepath.Abs(filepath.Join("somewhere", "demo"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
