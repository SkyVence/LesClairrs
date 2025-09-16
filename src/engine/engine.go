package engine

type Game interface {
	Init() Msg
	Update(Msg) (Model, Cmd)
	View() string
}

type engineModel struct {
	game Game
}

func Wrap(g Game) Model {
	return &engineModel{game: g}
}

func (e *engineModel) Init() Msg {
	return e.game.Init()
}

func (e *engineModel) Update(msg Msg) (Model, Cmd) {
	return e.game.Update(msg)
}

func (e *engineModel) View() string {
	return e.game.View()
}
