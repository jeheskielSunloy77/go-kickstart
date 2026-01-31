package scaffold

import "testing"

func TestToKebabCase(t *testing.T) {
	cases := map[string]string{
		"go-kickstart":   "go-kickstart",
		"Go Kickstart":   "go-kickstart",
		"GoKickstart":    "go-kickstart",
		"  My_App  ":     "my-app",
		"API V2 Service": "api-v2-service",
		"":               "",
	}

	for input, expected := range cases {
		if got := toKebabCase(input); got != expected {
			t.Fatalf("toKebabCase(%q) = %q, want %q", input, got, expected)
		}
	}
}
