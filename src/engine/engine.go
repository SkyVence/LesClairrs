package engine

type Game interface {
	Init() Msg
	Update(Msg) (Model, Cmd)
	View() string
}

type engineModel struct {
	game Game
}

// Wrap converts a Game into an engine-compatible Model
func Wrap(g Game) Model {
	return &engineModel{game: g}
}

// Init delegates to wrapped game's Init method
func (e *engineModel) Init() Msg {
	return e.game.Init()
}

// Update delegates to wrapped game's Update method
func (e *engineModel) Update(msg Msg) (Model, Cmd) {
	return e.game.Update(msg)
}

// View delegates to wrapped game's View method
func (e *engineModel) View() string {
	return e.game.View()
}
