package ui

import (
	"fmt"
	"strings"

	"projectred-rpg.com/term"
)

type Renderer struct{}

func (r Renderer) Render(ops []RenderOp) {
	term.Clear()

	for _, op := range ops {
		switch v := op.(type) {
		case TextOp:
			term.Move(v.Y, v.X)
			open := ""
			if v.Bold {
				open += term.Bold
			}
			if v.Fg != "" {
				open += term.Color(v.Fg)
			}
			fmt.Print(open + v.Text + term.Reset())
		case FillOp:
			for row := 0; row < v.H; row++ {
				term.Move(v.Y+row, v.X)
				fmt.Print(strings.Repeat(string(v.Ch), v.W))
			}
		case BorderOp:
			drawBorder(v)
		}
	}
}

func drawBorder(b BorderOp) {
	// horizontal
	for x := 0; x < b.W; x++ {
		uiPut(b.X+x, b.Y, b.Style.H)
		uiPut(b.X+x, b.Y+b.H-1, b.Style.H)
	}
	// vertical
	for y := 0; y < b.H; y++ {
		uiPut(b.X, b.Y+y, b.Style.V)
		uiPut(b.X+b.W-1, b.Y+y, b.Style.V)
	}
	// corners
	uiPut(b.X, b.Y, b.Style.TL)
	uiPut(b.X+b.W-1, b.Y, b.Style.TR)
	uiPut(b.X, b.Y+b.H-1, b.Style.BL)
	uiPut(b.X+b.W-1, b.Y+b.H-1, b.Style.BR)
}

func uiPut(x, y int, ch rune) {
	term.Move(y, x)
	fmt.Print(string(ch))
}
