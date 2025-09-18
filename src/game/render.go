package game

import (
	"fmt"
	"strings"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/loaders"
	"projectred-rpg.com/game/types"
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
	if err := loaders.LoadWorlds(); err != nil {
		// Handle error gracefully - could log it or show an error message
		// For now, we'll continue with empty cache
	}

	// Set up localization
	locManager := engine.GetLocalizationManager()
	locManager.SetLanguage("fr")

	// Later, to change language for all components
	//registry := ui.GetComponentRegistry()
	//registry.ChangeLanguage("en")

	// Create menu options
	menuOptions := []ui.MenuOption{
		{Label: locManager.Text("ui.menu.start"), Value: "start"},
		{Label: locManager.Text("ui.menu.settings"), Value: "settings"},
		{Label: locManager.Text("ui.menu.quit"), Value: "quit"},
	}
	// Load ASCII art from file
	menu, err := ui.NewMenuWithArtFromFile(locManager.Text("ui.menu.mainmenu"), menuOptions, "assets/logo.txt")
	if err != nil {
		// If loading the ASCII art fails, fall back to a simple menu without art
		menu = ui.NewMenu("", menuOptions)
	}
	return &model{
		state: stateMenu,
		menu:  menu,
		game: NewGameInstance(types.Class{
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
	// Route to specific state handlers for better organization
	switch m.state {
	case stateMenu:
		return m.renderMenuState()
	case stateGame:
		return m.renderGameState()
	case stateTransition:
		return m.renderTransitionState()
	default:
		return "Unknown state"
	}
}

// Separate rendering methods for better organization and maintainability

// renderMenuState handles menu state rendering
func (m *model) renderMenuState() string {
	return m.menu.View()
}

// renderGameState handles main game rendering with HUD and game world
func (m *model) renderGameState() string {
	// Update HUD with current player stats and location
	m.updateHUDStats()

	// Render game world content
	gameContent := m.gameSpace.RenderGameWorld(m.game.Player)

	// Combine HUD with game content
	return m.hud.RenderWithContent(gameContent)
}

// renderTransitionState handles loading/transition screen rendering
func (m *model) renderTransitionState() string {
	// Get next location information
	nextName, nameID, ok := m.game.PeekNext()

	// Create appropriate title
	title := m.createTransitionTitle(nextName, nameID, ok)

	// Update spinner animation
	_, cmd := m.spinner.Update(engine.TickNow())
	_ = cmd

	// Combine title with spinner
	content := title + "  " + m.spinner.View()
	return m.hud.RenderWithContent(content)
}

// updateHUDStats updates HUD with current player stats and location info
func (m *model) updateHUDStats() {
	player := m.game.Player
	worldID := m.game.CurrentWorld.WorldID
	stageID := m.game.CurrentStage.StageNb

	// Set player statistics
	m.hud.SetPlayerStats(
		player.Stats.CurrentHP,
		player.Stats.MaxHP,
		player.Stats.Level,
		int(player.Stats.Exp),
		player.Stats.NextLevelExp,
		worldID,
		stageID,
	)

	// Set location names from loaded world data
	if m.game.CurrentWorld != nil && m.game.CurrentStage != nil {
		m.hud.SetLocation(m.game.CurrentWorld.Name, m.game.CurrentStage.Name)
	}
}

// createTransitionTitle creates appropriate title for transition screens
func (m *model) createTransitionTitle(nextName string, nameID int, ok bool) string {
	locManager := engine.GetLocalizationManager()
	if !ok || nextName == "" {
		return locManager.Text("ui.menu.loading")
	}

	// Try to get translated name
	if translatedName := locManager.Text(fmt.Sprint("%s", "game.world"+fmt.Sprint(nameID)+".name")); translatedName != "" {
		return "Traveling to: " + translatedName
	}

	// Fallback to original name
	return "Traveling to: " + nextName
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

func (gr *GameRenderer) RenderGameWorld(player *types.Player) string {
	if gr.width <= 0 || gr.height <= 0 {
		return "Screen too small"
	}

	// Initialize game grid
	grid := gr.initializeGrid()

	// Render in organized layers
	gr.renderBackground(grid)
	gr.renderBorders(grid)
	gr.renderPlayer(grid, player)

	// Convert grid to string efficiently
	return gr.gridToString(grid)
}

// initializeGrid creates the base grid for rendering
func (gr *GameRenderer) initializeGrid() [][]rune {
	grid := make([][]rune, gr.height)
	for i := range grid {
		grid[i] = make([]rune, gr.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	return grid
}

// renderBackground creates background patterns and terrain
func (gr *GameRenderer) renderBackground(grid [][]rune) {
	// Simple background pattern - can be enhanced later
	for i := 1; i < gr.height-1; i++ {
		for j := 1; j < gr.width-1; j++ {
			grid[i][j] = ' '
		}
	}
}

// renderBorders draws the game area borders
func (gr *GameRenderer) renderBorders(grid [][]rune) {
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
}

// renderPlayer draws the player sprite on the grid
func (gr *GameRenderer) renderPlayer(grid [][]rune, player *types.Player) {
	spriteLines := strings.Split(player.GetSprite(), "\n")
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
}

// gridToString converts the grid to a string efficiently
func (gr *GameRenderer) gridToString(grid [][]rune) string {
	var builder strings.Builder
	builder.Grow(gr.width * gr.height) // Pre-allocate capacity

	for _, row := range grid {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}

	return strings.TrimRight(builder.String(), "\n")
}

// Extension methods for future systems - easy to implement when needed

// renderEnemies - placeholder for enemy rendering system
func (gr *GameRenderer) renderEnemies(grid [][]rune, enemies []types.Enemy) {
	// TODO: Implement when enemies have position data
	// This method provides a clear extension point for enemy rendering
}

// renderItems - placeholder for item rendering system
func (gr *GameRenderer) renderItems(grid [][]rune, items []interface{}) {
	// TODO: Implement when item system is added
	// This method provides a clear extension point for item rendering
}

// renderEffects - placeholder for visual effects system
func (gr *GameRenderer) renderEffects(grid [][]rune) {
	// TODO: Implement for particle effects, animations, etc.
	// This method provides a clear extension point for effects
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
