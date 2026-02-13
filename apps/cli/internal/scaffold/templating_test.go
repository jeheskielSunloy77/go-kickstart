package scaffold

import "testing"

func TestApplyTemplateConditions(t *testing.T) {
	input := "before\n{{IF_INCLUDE_WEB}}web-only {{PROJECT_NAME}}\n{{END_IF_INCLUDE_WEB}}after\n"

	withWeb := DefaultConfig()
	withWeb.IncludeWeb = true
	got := ApplyTemplateConditions(input, withWeb)
	if got != "before\nweb-only {{PROJECT_NAME}}\nafter\n" {
		t.Fatalf("unexpected content with web enabled: %q", got)
	}

	withoutWeb := DefaultConfig()
	withoutWeb.IncludeWeb = false
	got = ApplyTemplateConditions(input, withoutWeb)
	if got != "before\nafter\n" {
		t.Fatalf("unexpected content with web disabled: %q", got)
	}
}

func TestApplyTemplateConditionsThenReplaceTokens(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IncludeWeb = true

	input := "{{IF_INCLUDE_WEB}}hello {{PROJECT_NAME}}{{END_IF_INCLUDE_WEB}}"
	content := ApplyTemplateConditions(input, cfg)
	content = ReplaceTokens(content, map[string]string{"{{PROJECT_NAME}}": "demo"})

	if content != "hello demo" {
		t.Fatalf("unexpected rendered content: %q", content)
	}
}
