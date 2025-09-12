package components

import "projectred-rpg.com/ui"

type Label struct {
	Rect  ui.Rect
	Text  string
	Color string
	Bold  bool
}

func (l *Label) Init() []ui.Cmd             { return nil }
func (l *Label) SetBounds(r ui.Rect)        { l.Rect = r }
func (l *Label) Update(msg ui.Msg) []ui.Cmd { return nil }
func (l *Label) View() []ui.RenderOp {
	return []ui.RenderOp{
		ui.TextOp{
			X: l.Rect.X, Y: l.Rect.Y,
			Text: l.Text, Fg: l.Color, Bold: l.Bold,
		},
	}
}
