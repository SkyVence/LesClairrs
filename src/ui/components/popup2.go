package components

import (
	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/ui"
)

type Popup2 struct {
	Title   string
	Message string
	width   int
	height  int
}

func NewPopup2(title, message string) Popup2 {
	return Popup2{
		Title:   title,
		Message: message,
	}
}

func (p *Popup2) Update(msg ui.Msg) (Popup2, ui.Cmd) {
	switch msg := msg.(type) {
	case ui.SizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	}
	return *p, nil
}

func (p Popup2) View() string {
	return p.ViewWithSize(80, 24)
}

func (p Popup2) ViewWithSize(width, height int) string {
	// Style premier plan (plus visible)
	popup2BorderStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("196")). // Rouge vif
		Background(lipgloss.Color("52")).        // Rouge foncé
		Foreground(lipgloss.Color("255")).       // Blanc
		Padding(1, 2).
		Width(width - 15).  // Plus petit que popup1
		Height(height - 12) // Plus petit que popup1

	popup2TitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")). // Jaune
		Underline(true).
		Align(lipgloss.Center)

	popup2MessageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("52")).
		Align(lipgloss.Center)

	popup2FooterStyle := lipgloss.NewStyle().
		Faint(true).
		Bold(true).
		Align(lipgloss.Center)

	content := popup2TitleStyle.Render(p.Title) + "\n\n" +
		popup2MessageStyle.Render(p.Message) + "\n\n" +
		popup2FooterStyle.Render("[e] Fermer")

	return popup2BorderStyle.Render(content)
}
