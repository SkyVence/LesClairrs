package game

import (
	"fmt"
	"strings"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/ui"
)

type gameState int

const (
	stateMenu gameState = iota
	stateGame
	stateSettings
	stateTransition
)

type model struct {
	state     gameState
	menu      ui.Menu
	game      *Game
	gameSpace *GameRenderer
	hud       *ui.HUD
	spinner   ui.Spinner
	width     int
	height    int
}

func NewGame() *model {
	// Initialize all worlds at startup
	if err := LoadWorlds(); err != nil {
		// Handle error gracefully - could log it or show an error message
		// For now, we'll continue with empty cache

	}

	// Create menu options
	menuOptions := []ui.MenuOption{
		{Label: "Start Game", Value: "start"},
		{Label: "Settings", Value: "settings"},
		{Label: "Quit", Value: "quit"},
	}
	// Load ASCII art from file
	menu, err := ui.NewMenuWithArtFromFile("Game Options", menuOptions, "assets/logo.txt")
	if err != nil {
		// If loading the ASCII art fails, fall back to a simple menu without art
		menu = ui.NewMenu("", menuOptions)
	}
	return &model{
		state: stateMenu,
		menu:  menu,
		game: NewGameInstance(Class{
			Name:        "Cyber-Samurai",
			MaxHP:       100,
			Force:       1,
			Speed:       12,
			Defense:     8,
			Accuracy:    15,
			Description: "A swift and deadly warrior, excelling in close combat and agility.",
		}),
		hud:     ui.NewHud(),
		spinner: ui.NewSpinner(),
	}
}

func (m *model) Init() engine.Msg {
	return m.spinner.Init()
}

func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.menu, _ = m.menu.Update(msg)
		if m.gameSpace == nil {
			m.gameSpace = NewGameRenderer(msg.Width-1, msg.Height-m.hud.Height()-1)
		} else {
			m.gameSpace.UpdateSize(msg.Width-1, msg.Height-m.hud.Height()-1)
		}
		*m.hud, _ = m.hud.Update(msg)
		// Pass through size changes to spinner as well
		m.spinner, _ = m.spinner.Update(msg)

	case engine.KeyMsg:
		switch msg.Rune {
		case 'q':
			if m.state == stateGame {
				m.state = stateMenu
				return m, nil
			}
			return m, engine.Quit
		case 'n':
			// Demo: trigger transition to next stage/world
			if m.state == stateGame {
				if _, _, ok := m.game.PeekNext(); ok {
					m.state = stateTransition
					return m, m.spinner.Init()
				}
			}
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
		if m.state == stateTransition {
			// Advance spinner frames
			var cmd engine.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			// After a short time tick, perform the actual advance once
			if _, isTick := msg.(engine.TickMsg); isTick {
				// Move to next stage/world
				_ = m.game.Advance()
				m.state = stateGame
				return m, nil
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) View() string {
	lang, err := engine.Load("fr")
	if err != nil {
		return "Error loading language"
	}
	WorldID := m.game.CurrentWorld.WorldID
	StageId := m.game.CurrentStage.StageNb

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
			WorldID, // WorldID (placeholder)
			StageId, // StageID (placeholder)
		)

		// Set location names from the actual loaded world data
		if m.game.CurrentWorld != nil && m.game.CurrentStage != nil {
			m.hud.SetLocation(m.game.CurrentWorld.Name, m.game.CurrentStage.Name)
		}

		gameContent := m.gameSpace.RenderGameWorld(m.game.Player)
		return m.hud.RenderWithContent(gameContent)
	case stateTransition:
		// Show loading/transition overlay with spinner and next location name
		nextName, nameID, ok := m.game.PeekNext()
		if !ok {
			nextName = ""
		}
		title := "Loading"
		if nextName != "" {
			title = "Traveling to: " + lang.Text("level.world"+fmt.Sprint(nameID)+".name")
		}
		// force spinner to keep ticking
		_, cmd := m.spinner.Update(engine.TickNow())
		_ = cmd
		content := title + "  " + m.spinner.View()
		return m.hud.RenderWithContent(content)
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
