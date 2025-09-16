package aigen

// Backup of experimental game_renderer.go

import (
	"strings"
)

// GameRenderer handles rendering the game world with the stick man
type GameRenderer struct {
	width  int
	height int
}

// NewGameRenderer creates a new game renderer
func NewGameRenderer(width, height int) *GameRenderer {
	return &GameRenderer{
		width:  width,
		height: height,
	}
}

// RenderGameWorld creates a game screen with the stick man positioned correctly
func (gr *GameRenderer) RenderGameWorld(player *Player) string {
	// Safety check for invalid dimensions
	if gr.width <= 0 || gr.height <= 0 {
		return "Screen too small"
	}

	// Create a 2D grid to represent the screen
	lines := make([]string, gr.height)

	// Initialize with empty spaces
	for i := range lines {
		lines[i] = strings.Repeat(" ", gr.width)
	}

	// Get stick man sprite lines
	stickManLines := strings.Split(player.Sprite, "\n")

	// Position the stick man on the screen
	playerX, playerY := player.GetPosition()

	for i, line := range stickManLines {
		y := playerY + i
		if y >= 0 && y < gr.height && playerX >= 0 {
			// Replace characters in the line with the stick man sprite
			lineRunes := []rune(lines[y])
			for j, char := range line {
				x := playerX + j
				if x >= 0 && x < gr.width {
					lineRunes[x] = char
				}
			}
			lines[y] = string(lineRunes)
		}
	}

	// Add some basic world elements (optional borders)
	if gr.height > 0 && gr.width > 0 {
		// Top and bottom borders
		if gr.width > 0 {
			lines[0] = strings.Repeat("─", gr.width)
			if gr.height > 1 {
				lines[gr.height-1] = strings.Repeat("─", gr.width)
			}
		}

		// Side borders
		for i := 1; i < gr.height-1; i++ {
			if len(lines[i]) > 0 {
				lineRunes := []rune(lines[i])
				lineRunes[0] = '│'
				if len(lineRunes) > 1 {
					lineRunes[len(lineRunes)-1] = '│'
				}
				lines[i] = string(lineRunes)
			}
		}

		// Corners
		if gr.width > 0 && gr.height > 0 {
			lineRunes := []rune(lines[0])
			lineRunes[0] = '┌'
			if len(lineRunes) > 1 {
				lineRunes[len(lineRunes)-1] = '┐'
			}
			lines[0] = string(lineRunes)
		}

		if gr.width > 0 && gr.height > 1 {
			lineRunes := []rune(lines[gr.height-1])
			lineRunes[0] = '└'
			if len(lineRunes) > 1 {
				lineRunes[len(lineRunes)-1] = '┘'
			}
			lines[gr.height-1] = string(lineRunes)
		}
	}

	return strings.Join(lines, "\n")
}

// UpdateSize updates the renderer dimensions
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
