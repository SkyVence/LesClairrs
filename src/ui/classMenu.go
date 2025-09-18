package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
)

type ClassStats struct {
	MaxHP    int
	Force    int
	Speed    int
	Defense  int
	Accuracy int
}

type ClassMenuOption struct {
	Label string
	Value string
	Desc  string
	Stats ClassStats
}

type ClassMenu struct {
	Title    string
	Option   []ClassMenuOption
	Styles   ClassMenuStyles
	Loc      *engine.LocalizationManager
	selected int
	width    int
	height   int
}

type ClassMenuStyles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Normal      lipgloss.Style
	Description lipgloss.Style
	Stats       lipgloss.Style
	Sidebar     lipgloss.Style
}

func DefaultClassMenuStyles() ClassMenuStyles {
	return ClassMenuStyles{
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

func NewClassMenu(title string, options []ClassMenuOption, loc *engine.LocalizationManager, styles ...ClassMenuStyles) ClassMenu {
	menuStyles := DefaultClassMenuStyles()
	if len(styles) > 0 {
		menuStyles = styles[0]
	}

	return ClassMenu{
		Title:  title,
		Option: options,
		Styles: menuStyles,
		Loc:    loc,
	}
}

func (m ClassMenu) Update(msg engine.Msg) (ClassMenu, engine.Msg) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case engine.KeyMsg:
		switch msg.Rune {
		case '↓':
			if m.selected < len(m.Option)-1 {
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

// formatStats returns a compact, human-readable block of the starting stats.
func (m ClassMenu) formatStats(s ClassStats) string {
	// Labels are kept in English by default; can be localized later via keys if provided in assets
	return fmt.Sprintf(
		"%s: %d\n%s:  %d\n%s:  %d\n%s: %d\n%s: %d",
		m.localize("ui.class.menu.maxhp"), s.MaxHP,
		m.localize("ui.class.menu.force"), s.Force,
		m.localize("ui.class.menu.speed"), s.Speed,
		m.localize("ui.class.menu.defense"), s.Defense,
		m.localize("ui.class.menu.accuracy"), s.Accuracy,
	)
}

// localize returns a translated string if the input is a key present in the catalog;
// otherwise, it returns the input unchanged. It also falls back to the input if the
// translation is missing (to avoid showing ⟦key⟧).
func (m ClassMenu) localize(s string) string {
	if m.Loc == nil || s == "" {
		return s
	}
	tr := m.Loc.Text(s)
	if strings.HasPrefix(tr, "⟦") && strings.HasSuffix(tr, "⟧") {
		return s
	}
	return tr
}

// renderSidebar renders the description and stats for the selected option.
func (m ClassMenu) renderSidebar(opt ClassMenuOption, width int) string {
	if width <= 0 {
		return ""
	}
	descTitle := "ui.class.menu.description"
	statsTitle := "ui.class.menu.startingStats"

	descBlock := m.Styles.Description.
		Width(width).
		Render(m.localize(descTitle) + "\n\n" + m.localize(opt.Desc))

	statsBlock := m.Styles.Stats.
		Width(width).
		Render(m.localize(statsTitle) + "\n\n" + m.formatStats(opt.Stats))

	inner := lipgloss.JoinVertical(lipgloss.Left, descBlock, statsBlock)
	return m.Styles.Sidebar.Width(width).Render(inner)
}

func (m ClassMenu) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var menuItems []string

	menuItems = append(menuItems, m.Styles.Title.Render(m.localize(m.Title)))
	for i, option := range m.Option {
		var item string
		label := m.localize(option.Label)
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
	canTrySidebar := m.width >= minTotalForSidebar && len(m.Option) > 0

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
		rightContent := m.renderSidebar(m.Option[m.selected], rightW)
		right := lipgloss.Place(rightW, targetH, lipgloss.Left, lipgloss.Center, rightContent)
		content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left, gap, right)
	} else {
		// No room for sidebar; just keep the left menu centered by spacer
		content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left)
	}

	// Center the composed content vertically while keeping horizontal alignment left
	// (we already centered the menu horizontally via spacer)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Center,
		content,
	)
}

func (m ClassMenu) GetSelected() ClassMenuOption {
	if m.selected >= 0 && m.selected < len(m.Option) {
		return m.Option[m.selected]
	}
	return ClassMenuOption{}
}
