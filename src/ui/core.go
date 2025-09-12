package ui

type Rect struct {
	X, Y int
	W, H int
}

type Msg interface{}

type Component interface {
	Init() []Cmd
	Update(Msg) []Cmd
	View() []RenderOp
	SetBounds(Rect)
}

type Cmd func() Msg

type RenderOp interface {
	isOp()
}

type TextOp struct {
	X, Y int
	Fg   string
	Bg   string
	Bold bool
	Text string
}

func (TextOp) isOp() {}

type FillOp struct {
	X, Y, W, H int
	Ch         rune
}

func (FillOp) isOp() {}

type BorderOp struct {
	X, Y, W, H int
	Style      BorderStyle
}

func (BorderOp) isOp() {}

type BorderStyle struct {
	H, V           rune
	TL, TR, BL, BR rune
}
