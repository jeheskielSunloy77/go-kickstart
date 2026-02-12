package prompts

import "testing"

func TestDefaultModulePath(t *testing.T) {
	got := defaultModulePath("my-project")
	want := "github.com/yourorg/my-project"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
