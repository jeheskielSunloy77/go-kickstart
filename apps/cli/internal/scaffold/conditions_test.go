package scaffold

import "testing"

func TestShouldSkipForConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IncludeWeb = false
	cfg.IncludeDocker = false

	skip := ShouldSkipForConfig(cfg)

	if !skip("apps/web") {
		t.Fatalf("expected web directory to be skipped")
	}
	if !skip("apps/web/src/main.tsx") {
		t.Fatalf("expected web path to be skipped")
	}
	if !skip("packages/ui") {
		t.Fatalf("expected ui directory to be skipped")
	}
	if !skip("packages/ui/src/index.ts") {
		t.Fatalf("expected ui path to be skipped")
	}
	if !skip("docker-compose.yml") {
		t.Fatalf("expected docker compose to be skipped")
	}
	if skip("apps/api/main.go") {
		t.Fatalf("did not expect api path to be skipped")
	}
}
