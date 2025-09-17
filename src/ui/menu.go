package ui

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
	Class    string
}

// Message émis quand on valide une option
type MenuSelectMsg struct {
	Option MenuOption
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

func (m Menu) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
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
		case '\n', '\r':
			if len(m.Options) > 0 {
				opt := m.Options[m.selected]
				return m, func() engine.Msg { return MenuSelectMsg{Option: opt} }
			}
		}
	}
	return m, nil
}

func (m Menu) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, m.Styles.Title.Render(m.Title))
	lines = append(lines, "")

	for i, option := range m.Options {
		// Option "class" -> label dynamique
		label := option.Label
		if option.Value == "class" {
			label = "Class" // juste "Class" tout court
		}

		if i == m.selected {
			lines = append(lines, m.Styles.Selected.Render("▶ "+label))
		} else {
			lines = append(lines, m.Styles.Normal.Render("  "+label))
		}
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, lines...)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		menu,
	)
}

func (m Menu) Init() engine.Msg {
	return nil
}

func (m *Menu) SetClass(name string) {
	m.Class = name
}

func (m Menu) GetSelected() MenuOption {
	if m.selected >= 0 && m.selected < len(m.Options) {
		return m.Options[m.selected]
	}
	return MenuOption{}
}
