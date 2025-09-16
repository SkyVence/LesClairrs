package components

import (
	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/ui"
)

type Popup struct {
	Title   string
	Message string
	width   int
	height  int
}

func NewPopup(title, message string) Popup {
	return Popup{
		Title:   title,
		Message: message,
	}
}

func (p *Popup) Update(msg ui.Msg) (Popup, ui.Cmd) {
	switch msg := msg.(type) {
	case ui.SizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	}
	return *p, nil
}

func (p Popup) View() string {
	return p.ViewWithSize(80, 24)
}

func (p Popup) ViewWithSize(width, height int) string {
	// Style arrière-plan (plus subtil)
	popupBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")). // Gris plus clair
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("230")).
		Padding(1, 2).
		Width(width - 10). // Limite la largeur
		Height(height - 8) // Limite la hauteur

	popupTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center)

	popupMessageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("236")).
		Align(lipgloss.Center)

	popupFooterStyle := lipgloss.NewStyle().
		Faint(true).
		Align(lipgloss.Center)

	content := popupTitleStyle.Render(p.Title) + "\n\n" +
		popupMessageStyle.Render(p.Message) + "\n\n" +
		popupFooterStyle.Render("[q] Fermer")

	return popupBorderStyle.Render(content)
}
