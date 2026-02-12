package scaffold

import "testing"

func TestDefaultSkipEnvExample(t *testing.T) {
	cases := []struct {
		path string
		skip bool
	}{
		{path: ".env.production", skip: true},
		{path: ".env.example", skip: false},
		{path: ".env.example.tmpl", skip: false},
		{path: "apps/api/.env.example.tmpl", skip: false},
	}

	for _, c := range cases {
		if got := DefaultSkip(c.path); got != c.skip {
			t.Fatalf("DefaultSkip(%q) = %v, want %v", c.path, got, c.skip)
		}
	}
}
