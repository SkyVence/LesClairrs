package components

import (
	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
)

type MenuOption struct {
	Label string
	Value string
}

type MenuStyles struct {
	Title    lipgloss.Style
	Selected lipgloss.Style
	Normal   lipgloss.Style
}

func DefaultMenuStyles() MenuStyles {
	return MenuStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1),
		Selected: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8")).
			Background(lipgloss.Color("#654EA3")).
			Padding(0, 1),
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1),
	}
}

type Menu struct {
	Title    string
	Options  []MenuOption
	Styles   MenuStyles
	selected int
	width    int
	height   int
}

func NewMenu(title string, options []MenuOption, styles ...MenuStyles) Menu {
	menuStyles := DefaultMenuStyles()
	if len(styles) > 0 {
		menuStyles = styles[0]
	}

	return Menu{
		Title:   title,
		Options: options,
		Styles:  menuStyles,
	}
}

func (m Menu) Update(msg engine.Msg) (Menu, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case engine.KeyMsg:
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

	// Build menu content
	var menuItems []string

	// Add title
	menuItems = append(menuItems, m.Styles.Title.Render(m.Title))
	menuItems = append(menuItems, "") // Empty line

	// Add options
	for i, option := range m.Options {
		var item string
		if i == m.selected {
			item = m.Styles.Selected.Render("▶ " + option.Label)
		} else {
			item = m.Styles.Normal.Render("  " + option.Label)
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
