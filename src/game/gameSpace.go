package game

import (
	"strings"

	"projectred-rpg.com/game/types"
)

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
