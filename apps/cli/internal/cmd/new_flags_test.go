package cmd

import "testing"

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
