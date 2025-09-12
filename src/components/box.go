package components

import (
	"projectred-rpg.com/ui"
)

// Box represents a bordered container with an optional title and child component
type Box struct {
	Rect    ui.Rect
	Title   string
	Style   ui.BorderStyle
	Padding int
	Child   ui.Component
}

// DefaultBoxStyle defines the default border characters for the box
var DefaultBoxStyle = ui.BorderStyle{
	H: '─', V: '│', TL: '┌', TR: '┐', BL: '└', BR: '┘',
}

// Init initializes the box and its child component
func (b *Box) Init() []ui.Cmd {
	if b.Child != nil {
		return b.Child.Init()
	}
	return nil
}

// SetBounds sets the box boundaries and adjusts child component
func (b *Box) SetBounds(r ui.Rect) {
	b.Rect = r
	if b.Child != nil {
		child := ui.Rect{
			X: r.X + 1 + b.Padding,
			Y: r.Y + 1 + b.Padding,
			W: r.W - 2 - 2*b.Padding,
			H: r.H - 2 - 2*b.Padding,
		}
		b.Child.SetBounds(child)
	}
}

// Update handles input messages and updates box state
func (b *Box) Update(msg ui.Msg) []ui.Cmd {
	if b.Child != nil {
		return b.Child.Update(msg)
	}
	return nil
}

// View returns the render operations for the box
func (b *Box) View() []ui.RenderOp {
	ops := []ui.RenderOp{
		ui.BorderOp{X: b.Rect.X, Y: b.Rect.Y, W: b.Rect.W, H: b.Rect.H, Style: b.Style},
	}
	if b.Title != "" {
		ops = append(ops, ui.TextOp{
			X: b.Rect.X + 2, Y: b.Rect.Y,
			Text: " " + b.Title + " ", Bold: true,
		})
	}
	if b.Child != nil {
		ops = append(ops, b.Child.View()...)
	}
	return ops
}
