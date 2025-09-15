// main.go
package main

import (
	"log"

	"projectred-rpg.com/ui"

	"github.com/charmbracelet/lipgloss"
)

// model now holds our player's animation.
type model struct {
	player      ui.Animation
	playerWidth int
}

var _ ui.Model = (*model)(nil)

// newModel now loads the animation from the file.
func newModel() *model {
	// Load the animation frames from our .anim file.
	frames, err := ui.LoadAnimationFile("assets/animations/player-running.anim")
	if err != nil {
		// If the file can't be loaded, we can't run the game.
		log.Fatalf("Could not load animation file: %v", err)
	}

	maxWidth := 0
	for _, frame := range frames {
		w := lipgloss.Width(frame)
		if w > maxWidth {
			maxWidth = w
		}
	}

	return &model{
		player:      ui.NewAnimation(frames),
		playerWidth: maxWidth,
	}
}

// Init kicks off the animation for our component.
func (m *model) Init() ui.Msg {
	// We need to return the initial command from our animation component.
	return m.player.Init()()
}

// Update routes messages to the appropriate component.
func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
	var cmd ui.Cmd

	switch msg := msg.(type) {
	case ui.KeyMsg:
		switch msg.Rune {
		case 'q':
			return m, ui.Quit
		}

	// Any other message type (like TickMsg) is delegated to the player.
	default:
		m.player, cmd = m.player.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View composes the UI.
func (m *model) View() string {

	// Use lipgloss to place the animation and status text side-by-side.
	// Use ViewAligned instead of View to ensure consistent positioning
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.player.ViewAligned(m.playerWidth),
		"   ", // Some space
	)
}

func main() {
	p := ui.NewProgram(newModel(), ui.WithAltScreen())

	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
