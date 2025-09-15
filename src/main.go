package main

import (
	"log"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/ui"
	"projectred-rpg.com/ui/components"
)

type gameState int

const (
	stateMenu gameState = iota
	stateGame
	stateSettings
)

type model struct {
	state  gameState
	menu   components.Menu
	player ui.Animation
	hud    ui.HUD
	width  int
	height int
}

func newModel() *model {
	// Create menu options
	menuOptions := []components.MenuOption{
		{Label: "Start Game", Value: "start"},
		{Label: "Settings", Value: "settings"},
		{Label: "Quit", Value: "quit"},
	}

	menu := components.NewMenu("ProjectRed: RPG", menuOptions)

	// Load player animation
	frames, err := ui.LoadAnimationFile("assets/animations/player-running.anim")
	if err != nil {
		log.Fatalf("Could not load animation file: %v", err)
	}

	return &model{
		state:  stateMenu,
		menu:   menu,
		player: ui.NewAnimation(frames),
	}
}

func (m *model) Init() ui.Msg {
	if m.state == stateGame {
		return m.player.Init()()
	}
	return nil
}

func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
	switch msg := msg.(type) {
	case ui.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu, _ = m.menu.Update(msg)

	case ui.KeyMsg:
		switch msg.Rune {
		case 'q':
			if m.state == stateGame {
				m.state = stateMenu
				return m, nil
			}
			return m, ui.Quit
		case '\r', '\n', ' ': // Enter key
			if m.state == stateMenu {
				selected := m.menu.GetSelected()
				switch selected.Value {
				case "start":
					m.state = stateGame
					return m, m.player.Init()
				case "quit":
					return m, ui.Quit
				}
			}
		}

		// Route to appropriate component
		if m.state == stateMenu {
			m.menu, _ = m.menu.Update(msg)
		}

	default:
		// Route other messages (like TickMsg) to game components
		if m.state == stateGame {
			var cmd ui.Cmd
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
		gameContent := m.player.View()
		hudOverlay := m.hud.View()

		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Left,
			gameContent,
		) + "\n" + hudOverlay
	default:
		return "Unknown state"
	}
}

func main() {
	p := ui.NewProgram(newModel(), ui.WithAltScreen())
	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
