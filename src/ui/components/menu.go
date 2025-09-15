package components

import (
	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/ui"
)

type MenuOption struct {
	Label string
	Value string
}

type Menu struct {
	Title    string
	Options  []MenuOption
	selected int
	width    int
	height   int
}

func NewMenu(title string, options []MenuOption) Menu {
	return Menu{
		Title:   title,
		Options: options,
	}
}

func (m Menu) Update(msg ui.Msg) (Menu, ui.Cmd) {
	switch msg := msg.(type) {
	case ui.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ui.KeyMsg:
		switch msg.Rune {
		case '↓':
			if m.selected < len(m.Options)-1 {
				m.selected++
			}
		case '↑':
			if m.selected > 0 {
				m.selected--
			}
		}
	}
	return m, nil
}

func (m Menu) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#EE6FF8")).
		Background(lipgloss.Color("#654EA3")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1)

	// Build menu content
	var menuItems []string

	// Add title
	menuItems = append(menuItems, titleStyle.Render(m.Title))
	menuItems = append(menuItems, "") // Empty line

	// Add options
	for i, option := range m.Options {
		var item string
		if i == m.selected {
			item = selectedStyle.Render("▶ " + option.Label)
		} else {
			item = normalStyle.Render("  " + option.Label)
		}
		menuItems = append(menuItems, item)
	}

	// Join menu items
	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Center the menu on screen
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		menu,
	)
}

func (m Menu) GetSelected() MenuOption {
	if m.selected >= 0 && m.selected < len(m.Options) {
		return m.Options[m.selected]
	}
	return MenuOption{}
}
