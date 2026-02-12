package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintSummary(path string) {
	header := SectionTitleStyle().Render("✅ Repo spawned successfully")
	pathLine := lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true).Render(path)
	steps := []string{
		"1) `cd` into the project",
		"2) review generated `.env` files",
		"3) run the project wide install command `bun install`",
		"4) run the dev servers with `bun run dev` and read the README for more info",
	}

	fmt.Printf("\n%s\n", header)
	fmt.Printf("📦 Destination: %s\n", pathLine)
	fmt.Println(HintStyle().Render("Next moves:"))
	for _, step := range steps {
		fmt.Printf("   %s\n", step)
	}
	fmt.Println(HintStyle().Render("GG. Build cool things."))
}
