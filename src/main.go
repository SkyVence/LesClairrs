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
	statePopup
	statePopupDouble
)

type model struct {
	state       gameState
	menu        components.Menu
	player      ui.Animation
	popup       *components.Popup
	popup2      *components.Popup2
	playerWidth int
	width       int
	height      int
}

func newModel() *model {
	// Menu
	menuOptions := []components.MenuOption{
		{Label: "Start Game", Value: "start"},
		{Label: "Settings", Value: "settings"},
		{Label: "Quit", Value: "quit"},
	}
	menu := components.NewMenu("🎮 My Game", menuOptions)

	// Player
	frames, err := ui.LoadAnimationFile("assets/animations/player-running.anim")
	if err != nil {
		log.Fatalf("Could not load animation file: %v", err)
	}

	// Popups
	popup := components.NewPopup("Paramètres", "Réglez vos paramètres ici.")
	popup2 := components.NewPopup2("Confirmation", "Êtes-vous sûr ?")

	return &model{
		state:  stateMenu,
		menu:   menu,
		player: ui.NewAnimation(frames),
		popup:  &popup,
		popup2: &popup2,
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

		// Donner les dimensions aux popups
		if m.state == statePopup {
			*m.popup, _ = m.popup.Update(msg)
		}
		if m.state == statePopupDouble {
			*m.popup, _ = m.popup.Update(msg)
			*m.popup2, _ = m.popup2.Update(msg)
		}

	case ui.KeyMsg:
		switch msg.Rune {
		case 'q':
			if m.state == stateGame {
				m.state = stateMenu
				return m, nil
			}
			if m.state == statePopup || m.state == statePopupDouble {
				m.state = stateMenu
				return m, nil
			}
			return m, ui.Quit
		case 'e':
			if m.state == statePopupDouble {
				m.state = stateMenu
				return m, nil
			}
		case '\r', '\n', ' ': // Enter
			if m.state == stateMenu {
				selected := m.menu.GetSelected()
				switch selected.Value {
				case "start":
					m.state = stateGame
					return m, m.player.Init()
				case "settings":
					m.state = statePopupDouble // Ouvre les deux popups
					return m, nil
				case "quit":
					return m, ui.Quit
				}
			}
		}

		// Donner les touches aux composants
		if m.state == stateMenu {
			m.menu, _ = m.menu.Update(msg)
		}
		if m.state == statePopup {
			*m.popup, _ = m.popup.Update(msg)
		}
		if m.state == statePopupDouble {
			*m.popup, _ = m.popup.Update(msg)
			*m.popup2, _ = m.popup2.Update(msg)
		}

	default:
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
		return m.player.View()
	case statePopup:
		return m.popup.ViewWithSize(m.width, m.height)
	case statePopupDouble:
		// ============ DEUX POPUPS SIMULTANÉS AVEC SUPERPOSITION ============
		
		// Popup1 (arrière-plan) - position haut-gauche du centre
		popup1Content := m.popup.ViewWithSize(35, 12)
		popup1Positioned := lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().
				MarginTop(-4).     // Vers le haut
				MarginLeft(-8).    // Vers la gauche
				Render(popup1Content),
		)
		
		// Popup2 (premier plan) - position bas-droite du centre (déborde sur popup1)
		popup2Content := m.popup2.ViewWithSize(35, 12)
		popup2Positioned := lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().
				MarginTop(2).      // Vers le bas
				MarginLeft(6).     // Vers la droite
				Render(popup2Content),
		)
		
		// Technique de superposition avec JoinHorizontal pour garder les deux visibles
		leftSide := popup1Positioned
		rightSide := popup2Positioned
		
		// Combiner les deux côtés en gardant la superposition
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			leftSide,
			rightSide,
		)

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
