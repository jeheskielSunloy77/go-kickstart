package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle  = lipgloss.NewStyle().Bold(true)
	MutedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	AccentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
)
