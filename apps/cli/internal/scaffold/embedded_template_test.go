package scaffold

import (
	"io/fs"
	"testing"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/templates"
)

func TestEmbeddedMonorepoContainsAPIApp(t *testing.T) {
	sub, err := fs.Sub(templates.MonorepoFS, "monorepo")
	if err != nil {
		t.Fatalf("sub fs: %v", err)
	}
	info, err := fs.Stat(sub, "apps/api")
	if err != nil {
		t.Fatalf("expected apps/api in embedded templates: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected apps/api to be a directory")
	}
}
