package scaffold

import (
	"strings"
	"testing"
)

func TestMergeEnvExample(t *testing.T) {
	input := "FOO=bar\n# Comment\nBAZ=qux\n"
	overrides := map[string]string{
		"FOO": "override",
		"NEW": "value",
	}
	got := MergeEnvExample(input, overrides)
	if !containsLine(got, "FOO=override") {
		t.Fatalf("expected override to be applied")
	}
	if !containsLine(got, "NEW=value") {
		t.Fatalf("expected new key to be appended")
	}
}

func containsLine(input, line string) bool {
	for _, l := range strings.Split(input, "\n") {
		if l == line {
			return true
		}
	}
	return false
}
