package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func HuhTheme() *huh.Theme {
	t := huh.ThemeCharm()

	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color("51")).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color("114"))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color("205")).Bold(true)
	t.Focused.Option = t.Focused.Option.Foreground(lipgloss.Color("252"))
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color("16")).Background(lipgloss.Color("45")).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color("207")).Bold(true)
	t.Focused.Next = t.Focused.Next.Foreground(lipgloss.Color("220")).Bold(true)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(lipgloss.Color("205"))
	t.Blurred.Title = t.Blurred.Title.Foreground(lipgloss.Color("244"))

	return t
}

func BannerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("51"))
}

func SectionTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true)
}

func HintStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("114"))
}
