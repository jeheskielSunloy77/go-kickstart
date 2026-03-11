package cmd

import (
	"path/filepath"
	"testing"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func TestConfigFromFlagsRequiresModule(t *testing.T) {
	_, err := configFromFlags([]string{"demo"}, newFlags{})
	if err == nil {
		t.Fatalf("expected error when module is missing")
	}
}

func TestConfigFromFlagsUsesArgs(t *testing.T) {
	flags := newFlags{modulePath: "github.com/acme/demo"}
	cfg, err := configFromFlags([]string{"demo", "./out"}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ProjectName != "demo" {
		t.Fatalf("expected project name from args")
	}
	if cfg.Destination == "" {
		t.Fatalf("expected destination")
	}
}

func TestConfigFromFlagsDefaultsToProjectDirectoryInCwd(t *testing.T) {
	flags := newFlags{modulePath: "github.com/acme/demo"}
	cfg, err := configFromFlags([]string{"demo"}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, err := filepath.Abs(filepath.Join(".", "demo"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if cfg.Destination != want {
		t.Fatalf("expected %s, got %s", want, cfg.Destination)
	}
}

func TestConfigFromFlagsUsesRelativeBasePath(t *testing.T) {
	flags := newFlags{modulePath: "github.com/acme/demo"}
	cfg, err := configFromFlags([]string{"demo", "some-dir"}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, err := filepath.Abs(filepath.Join("some-dir", "demo"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if cfg.Destination != want {
		t.Fatalf("expected %s, got %s", want, cfg.Destination)
	}
}

func TestConfigFromFlagsUsesAbsoluteBasePath(t *testing.T) {
	base := t.TempDir()
	flags := newFlags{modulePath: "github.com/acme/demo"}
	cfg, err := configFromFlags([]string{"demo", base}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(base, "demo")
	if cfg.Destination != want {
		t.Fatalf("expected %s, got %s", want, cfg.Destination)
	}
}

func TestConfigFromFlagsSmartDedupeWhenBaseEndsWithProjectName(t *testing.T) {
	parent := t.TempDir()
	base := filepath.Join(parent, "demo")
	flags := newFlags{modulePath: "github.com/acme/demo"}
	cfg, err := configFromFlags([]string{"demo", base}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Destination != base {
		t.Fatalf("expected %s, got %s", base, cfg.Destination)
	}
}

func TestConfigFromFlagsParsesObservability(t *testing.T) {
	flags := newFlags{
		modulePath:    "github.com/acme/demo",
		observability: "grafana-oss",
	}
	cfg, err := configFromFlags([]string{"demo"}, flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Observability != scaffold.ObservabilityGrafanaOSS {
		t.Fatalf("expected grafana-oss observability, got %s", cfg.Observability)
	}
}

func TestConfigFromFlagsRejectsUnknownObservability(t *testing.T) {
	flags := newFlags{
		modulePath:    "github.com/acme/demo",
		observability: "bogus",
	}
	if _, err := configFromFlags([]string{"demo"}, flags); err == nil {
		t.Fatalf("expected error for unsupported observability")
	}
}
