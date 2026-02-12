package ui

import (
	"strings"
	"testing"
)

func TestWelcomeCardIncludesBannerMetadata(t *testing.T) {
	version := "0.1.0"
	repo := "https://github.com/jeheskielSunloy77/go-kickstart"

	card := WelcomeCard(version, repo)

	if !strings.Contains(card, "go-kickstart") {
		t.Fatalf("welcome card should include go-kickstart banner text")
	}
	if !strings.Contains(card, "Version: "+version) {
		t.Fatalf("welcome card should include version")
	}
	if !strings.Contains(card, "Repo: "+repo) {
		t.Fatalf("welcome card should include repo URL")
	}
	if !strings.Contains(card, ContributionLine()) {
		t.Fatalf("welcome card should include contribution line")
	}
}
