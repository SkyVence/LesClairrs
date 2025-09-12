package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/ui"
)

type model struct {
	count   int
	spinner ui.Spinner
}

var _ ui.Model = &model{}

func newModel() *model {
	return &model{
		count:   0,
		spinner: ui.NewSpinner(),
	}
}

func (m *model) Init() ui.Msg {
	// Return nil and let the spinner start in the first Update call
	return nil
}

func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
	var cmd ui.Cmd
	var spinnerCmd ui.Cmd

	// Always update the spinner with any message
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case ui.KeyMsg:
		switch msg.Rune {
		case 'q':
			return m, ui.Quit
		case 'j', 'd':
			m.count--
		case 'k', 'u':
			m.count++
		}
	case ui.TickMsg:
		// Spinner already updated above
		return m, spinnerCmd
	default:
		// For the first nil message, start the spinner
		if msg == nil {
			return m, m.spinner.Init()
		}
	}

	// Return spinner command if we have one, otherwise nil
	if spinnerCmd != nil {
		return m, spinnerCmd
	}

	return m, cmd
}

func (m *model) View() string {
	// Render the counter part of the view.
	counterView := fmt.Sprintf("Count: %d", m.count)

	// Render the help text.
	helpView := "[j/d] to decrement | [k/u] to increment | [q] to quit"

	// Combine the spinner and counter on the same line.
	mainView := lipgloss.JoinHorizontal(lipgloss.Top,
		m.spinner.View()+" ", // Get the spinner's view and add a space
		counterView,
	)

	// Stack the main view and help text vertically.
	return lipgloss.JoinVertical(lipgloss.Left,
		mainView,
		"", // Add a blank line
		helpView,
	)
}

func main() {
	p := ui.NewProgram(newModel(), ui.WithAltScreen())

	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
