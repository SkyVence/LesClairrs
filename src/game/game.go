package game

import (
	"log"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/ui"
)

type gameState int

const (
	stateMenu gameState = iota
	stateGame
	stateSettings
)

type model struct {
	state  gameState
	menu   ui.Menu
	player engine.Animation
	hud    *ui.HUD
	width  int
	height int
}

func NewGame() *model {
	// Create menu options
	menuOptions := []ui.MenuOption{
		{Label: "Start Game", Value: "start"},
		{Label: "Settings", Value: "settings"},
		{Label: "Quit", Value: "quit"},
	}

	menu := ui.NewMenu("ProjectRed: RPG", menuOptions)

	// Load player animation
	frames, err := engine.LoadAnimationFile("assets/animations/loader.anim")
	if err != nil {
		log.Fatalf("Could not load animation file: %v", err)
	}

	return &model{
		state:  stateMenu,
		menu:   menu,
		player: engine.NewAnimation(frames),
		hud:    ui.NewHud(),
	}
}

func (m *model) Init() engine.Msg {
	if m.state == stateGame {
		return m.player.Init()()
	}
	return nil
}

func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu, _ = m.menu.Update(msg)
		*m.hud, _ = m.hud.Update(msg)

	case engine.KeyMsg:
		switch msg.Rune {
		case 'q':
			if m.state == stateGame {
				m.state = stateMenu
				return m, nil
			}
			return m, engine.Quit
		case '\r', '\n', ' ': // Enter key
			if m.state == stateMenu {
				selected := m.menu.GetSelected()
				switch selected.Value {
				case "start":
					m.state = stateGame
					return m, m.player.Init()
				case "quit":
					return m, engine.Quit
				}
			}
		case '↑', '↓', '←', '→':
			_ = msg.Rune
		}

		if m.state == stateMenu {
			m.menu, _ = m.menu.Update(msg)
		}

	default:
		if m.state == stateGame {
			var cmd engine.Cmd
			m.player, cmd = m.player.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) View() string {
	switch m.state {
	case stateMenu:
		return m.menu.View()
	case stateGame:
		m.hud.SetPlayerStats(100, 100, 2, 75, 200, "Cyber District")
		gameContent := m.player.View()
		return m.hud.RenderWithContent(gameContent)
	default:
		return "Unknown state"
	}
}
