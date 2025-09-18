package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
)

type SettingsMenuOption struct {
	Label string
	Value string
}

type SettingsMenu struct {
	Title    string
	Options  []SettingsMenuOption
	Styles   SettingsMenuStyles
	Loc      *engine.LocalizationManager
	selected int
	width    int
	height   int
}

type SettingsMenuStyles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Normal      lipgloss.Style
	Description lipgloss.Style
	Stats       lipgloss.Style
	Sidebar     lipgloss.Style
}

func DefaultSettingsMenuStyles() SettingsMenuStyles {
	return SettingsMenuStyles{
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
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C0C0C0")).
			Padding(1, 1).
			MarginTop(1),
		Stats: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(1, 1).
			MarginTop(1),
		Sidebar: lipgloss.NewStyle().
			Background(lipgloss.Color("#1F1F2E")).
			Padding(0, 1),
	}
}

func NewSettingsMenu(title string, options []SettingsMenuOption, loc *engine.LocalizationManager, styles ...SettingsMenuStyles) SettingsMenu {
	menuStyles := DefaultSettingsMenuStyles()
	if len(styles) > 0 {
		menuStyles = styles[0]
	}

	return SettingsMenu{
		Title:   title,
		Options: options,
		Styles:  menuStyles,
		Loc:     loc,
	}
}

func (m SettingsMenu) localize(s string) string {
	if m.Loc == nil || s == "" {
		return s
	}
	tr := m.Loc.Text(s)
	if strings.HasPrefix(tr, "⟦") && strings.HasSuffix(tr, "⟧") {
		return s
	}
	return tr
}

func (m SettingsMenu) Update(msg engine.Msg) (SettingsMenu, engine.Msg) {
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

func (m SettingsMenu) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var menuItems []string

	menuItems = append(menuItems, m.Styles.Title.Render(m.localize(m.Title)))
	for i, option := range m.Options {
		var item string
		label := m.localize("ui.settings.menu." + option.Label)
		if i == m.selected {
			item = m.Styles.Selected.Render("▶ " + label)
		} else {
			item = m.Styles.Normal.Render("  " + label)
		}
		menuItems = append(menuItems, item)
	}

	// Left column (menu)
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Try to show a sidebar without affecting the horizontal centering of the left menu
	// The left menu remains horizontally centered; the sidebar is placed to its right
	// if there is enough space.
	const minTotalForSidebar = 44 // rough minimum to keep things readable
	canTrySidebar := m.width >= minTotalForSidebar && len(m.Options) > 0

	// Base widths and gap
	gapW := 2
	leftW := m.width * 2 / 5
	if leftW < 18 {
		leftW = 18
	}
	// Ensure the left column cannot exceed the available width
	if leftW > m.width {
		leftW = m.width
	}

	// Keep left column horizontally centered irrespective of the sidebar
	// We'll compute a left spacer (margin) so that the left column is centered.
	leftMargin := 0
	if m.width > leftW {
		leftMargin = (m.width - leftW) / 2
	}

	// Decide if we can fit a sidebar to the RIGHT of the centered left column
	rightW := 0
	if canTrySidebar {
		// Space remaining on the right side of the centered left column
		availableRight := m.width - (leftMargin + leftW) - gapW
		if availableRight >= 20 {
			rightW = availableRight
		}
	}

	// Build fixed-height boxes so total height stays constant and can be centered
	targetH := m.height
	if targetH < 1 {
		targetH = 1
	}

	// Left column box (content centered vertically inside its box)
	left := lipgloss.Place(leftW, targetH, lipgloss.Left, lipgloss.Center, leftColumn)

	// Horizontal spacers with fixed height to match boxes
	spacer := lipgloss.NewStyle().Width(leftMargin).Height(targetH).Render("")
	gap := lipgloss.NewStyle().Width(gapW).Height(targetH).Render("")

	var content string
	if rightW > 0 {
		rightContent := m.renderSidebar(m.Options[m.selected], rightW)
		right := lipgloss.Place(rightW, targetH, lipgloss.Left, lipgloss.Center, rightContent)
		content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left, gap, right)
	} else {
		// No room for sidebar; just keep the left menu centered by spacer
		content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left)
	}

	// Center the composed content vertically while keeping horizontal alignment left
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Center,
		content,
	)
}

// renderSidebar renders the description and stats for the selected option.
func (m SettingsMenu) renderSidebar(opt SettingsMenuOption, width int) string {
	if width <= 0 {
		return ""
	}
	descTitle := "ui.settings.menu.exampleSidebarTitle"
	statsTitle := "ui.settings.menu.exampleSidebarDesc"
	langTitle := "ui.settings.menu.language"

	descBlock := m.Styles.Description.
		Width(width).
		Render(m.localize(descTitle))

	statsBlock := m.Styles.Stats.
		Width(width).
		Render(m.localize(statsTitle))

	langBlock := m.Styles.Stats.
		Width(width).
		Render(m.localize(langTitle) + ": " + m.localize(opt.Value))

	inner := lipgloss.JoinVertical(lipgloss.Left, descBlock, statsBlock, langBlock)
	return m.Styles.Sidebar.Width(width).Render(inner)
}

func (m SettingsMenu) GetSelected() SettingsMenuOption {
	if m.selected >= 0 && m.selected < len(m.Options) {
		return m.Options[m.selected]
	}
	return SettingsMenuOption{}
}
