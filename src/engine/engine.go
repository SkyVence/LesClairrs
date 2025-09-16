package engine

import "projectred-rpg.com/ui"

type Game interface {
	Init() ui.Msg
	Update(ui.Msg) (ui.Model, ui.Cmd)
	View() string
}

type engineModel struct {
	game Game
}

func Wrap(g Game) ui.Model {
    return &engineModel{game: g}
}

func (e *engineModel) Init() ui.Msg {
    return e.game.Init()
}

func (e *engineModel) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    return e.game.Update(msg)
}

func (e *engineModel) View() string {
    return e.game.View()
}