package validate

import "testing"

func TestProjectName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"ok", "myapp", false},
		{"ok with n", "admin", false},
		{"slash", "my/app", true},
		{"backslash", "my\\app", true},
		{"newline", "my\napp", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ProjectName(c.input)
			if c.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestModulePath(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"ok", "github.com/acme/foo", false},
		{"no slash", "github.com", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ModulePath(c.input)
			if c.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
