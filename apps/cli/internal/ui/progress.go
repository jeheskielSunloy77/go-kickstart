package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type spinnerModel struct {
	spinner spinner.Model
	message string
	err     error
	done    <-chan error
}

type spinnerDoneMsg struct{ err error }

type spinnerTickMsg struct{}

func RunWithSpinner(message string, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	m := spinnerModel{
		spinner: spinner.New(),
		message: message,
		done:    done,
	}

	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return err
	}
	if sm, ok := final.(spinnerModel); ok {
		return sm.err
	}
	return nil
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForSpinner(m.done))
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerDoneMsg:
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.err = fmt.Errorf("cancelled")
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

func waitForSpinner(done <-chan error) tea.Cmd {
	return func() tea.Msg {
		return spinnerDoneMsg{err: <-done}
	}
}
