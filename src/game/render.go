package game

import (
	"strings"

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
	state     gameState
	menu      ui.Menu
	game      *Game
	gameSpace *GameRenderer
	hud       *ui.HUD
	width     int
	height    int
}

func NewGame() *model {
	// Create menu options
	menuOptions := []ui.MenuOption{
		{Label: "Start Game", Value: "start"},
		{Label: "Settings", Value: "settings"},
		{Label: "Quit", Value: "quit"},
	}

	menu := ui.NewMenu("ProjectRed: RPG", menuOptions)

	return &model{
		state: stateMenu,
		menu:  menu,
		game:  NewGameInstance(),
		hud:   ui.NewHud(),
	}
}

func (m *model) Init() engine.Msg {
	return nil
}

func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu, _ = m.menu.Update(msg)
		if m.gameSpace == nil {
			m.gameSpace = NewGameRenderer(msg.Width, msg.Height-m.hud.Height())
		} else {
			m.gameSpace.UpdateSize(msg.Width, msg.Height-m.hud.Height())
		}
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
					return m, nil
				case "quit":
					return m, engine.Quit
				}
			}
		case '↑', '↓', '←', '→':
			if m.state == stateGame {
				m.game.Player.Move(msg.Rune, m.gameSpace.width, m.gameSpace.height)
			}
		}

		if m.state == stateMenu {
			m.menu, _ = m.menu.Update(msg)
		}

	default:
		if m.state == stateGame {
			// var cmd engine.Cmd
			// m.gameSpace, cmd = m.gameSpace.Update(msg)
			// return m, cmd
		}
	}

	return m, nil
}

func (m *model) View() string {
	switch m.state {
	case stateMenu:
		return m.menu.View()
	case stateGame:
		player := m.game.Player
		m.hud.SetPlayerStats(
			player.Stats.CurrentHP,
			player.Stats.MaxHP,
			player.Stats.Level,
			int(player.Stats.Exp),
			player.Stats.NextLevelExp,
			"Cyber District",
		)
		gameContent := m.gameSpace.RenderGameWorld(m.game.Player)
		return m.hud.RenderWithContent(gameContent)
	default:
		return "Unknown state"
	}
}

// Game space renderer

type GameRenderer struct {
	width  int
	height int
}

func NewGameRenderer(width, height int) *GameRenderer {
	return &GameRenderer{
		width:  width,
		height: height,
	}
}

func (gr *GameRenderer) RenderGameWorld(player *Player) string {
	if gr.width <= 0 || gr.height <= 0 {
		return "Screen too small"
	}

	// Create a 2D grid for the game world
	grid := make([][]rune, gr.height)
	for i := range grid {
		grid[i] = make([]rune, gr.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Draw borders
	for i := 0; i < gr.height; i++ {
		if i == 0 || i == gr.height-1 {
			for j := 0; j < gr.width; j++ {
				switch {
				case (i == 0 && j == 0):
					grid[i][j] = '┌'
				case (i == 0 && j == gr.width-1):
					grid[i][j] = '┐'
				case (i == gr.height-1 && j == 0):
					grid[i][j] = '└'
				case (i == gr.height-1 && j == gr.width-1):
					grid[i][j] = '┘'
				case i == 0 || i == gr.height-1:
					grid[i][j] = '─'
				}
			}
		} else {
			grid[i][0] = '│'
			grid[i][gr.width-1] = '│'
		}
	}

	// Draw player sprite
	spriteLines := strings.Split(player.sprite, "\n")
	playerX, playerY := player.GetPosition()

	for i, line := range spriteLines {
		y := playerY + i
		if y >= 1 && y < gr.height-1 { // Ensure y is within borders
			for j, char := range line {
				x := playerX + j
				if x >= 1 && x < gr.width-1 { // Ensure x is within borders
					grid[y][x] = char
				}
			}
		}
	}

	// Convert grid to a single string
	var builder strings.Builder
	for _, row := range grid {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}

func (gr *GameRenderer) UpdateSize(width, height int) {
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}
	gr.width = width
	gr.height = height
}
