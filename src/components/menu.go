package components

import (
	"projectred-rpg.com/ui"
)

// Menu represents a vertical list of selectable items
type Menu struct {
	Rect  ui.Rect
	Items []string
	Index int
}

// MoveUp represents an upward movement command
type MoveUp struct{}

// MoveDown represents a downward movement command
type MoveDown struct{}

// Enter represents a selection command
type Enter struct{}

// Init initializes the menu component
func (m *Menu) Init() []ui.Cmd { return nil }

// SetBounds sets the menu's display boundaries
func (m *Menu) SetBounds(r ui.Rect) { m.Rect = r }

// Update handles input messages and updates menu state
func (m *Menu) Update(msg ui.Msg) []ui.Cmd {
	switch msg.(type) {
	case MoveUp:
		if m.Index > 0 {
			m.Index--
		}
	case MoveDown:
		if m.Index < len(m.Items)-1 {
			m.Index++
		}
	case Enter:
		// TODO: Implement selection handling
		return nil
	}
	return nil
}

// View returns the render operations for the menu
func (m *Menu) View() []ui.RenderOp {
	var ops []ui.RenderOp
	for i, it := range m.Items {
		color := ""
		bold := false
		prefix := "  "
		if i == m.Index {
			color = "36" // Cyan
			bold = true
			prefix = "> "
		}
		ops = append(ops, ui.TextOp{
			X: m.Rect.X, Y: m.Rect.Y + i,
			Text: prefix + it, Fg: color, Bold: bold,
		})
	}
	return ops
}
