package game

import (
	"strings"

	"projectred-rpg.com/game/types"
)

// Game space renderer with viewport/camera
type GameRenderer struct {
	width   int
	height  int
	tileMap *types.TileMap
	viewX   int // top-left map X of the viewport
	viewY   int // top-left map Y of the viewport
	// draw state for current frame
	scale int // integer scale (>=1); 1 = no scale, >1 = scaled up
	offX  int // interior offset X when scaled to center
	offY  int // interior offset Y when scaled to center
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

	// Adjust camera based on player and map
	gr.updateViewport(player)

	// Initialize draw state and grid
	gr.scale, gr.offX, gr.offY = 1, 0, 0
	grid := gr.initializeGrid()

	// Render in organized layers
	gr.renderBackground(grid)
	gr.renderMap(grid)
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
				case i == 0 && j == 0:
					grid[i][j] = '┌'
				case i == 0 && j == gr.width-1:
					grid[i][j] = '┐'
				case i == gr.height-1 && j == 0:
					grid[i][j] = '└'
				case i == gr.height-1 && j == gr.width-1:
					grid[i][j] = '┘'
				default:
					grid[i][j] = '─'
				}
			}
		} else {
			grid[i][0] = '│'
			grid[i][gr.width-1] = '│'
		}
	}
}

// SetMap sets the current tile map to render
func (gr *GameRenderer) SetMap(tm *types.TileMap) { gr.tileMap = tm }

// MapSize returns underlying map dimensions (0,0 if none)
func (gr *GameRenderer) MapSize() (int, int) {
	if gr.tileMap == nil {
		return 0, 0
	}
	return gr.tileMap.Width, gr.tileMap.Height
}

// updateViewport recenters or clamps the viewport based on player position and map size
func (gr *GameRenderer) updateViewport(player *types.Player) {
	if gr.tileMap == nil || player == nil {
		gr.viewX, gr.viewY = 0, 0
		return
	}
	visW := max(1, gr.width-2)
	visH := max(1, gr.height-2)
	mapW := max(1, gr.tileMap.Width)
	mapH := max(1, gr.tileMap.Height)

	px, py := player.GetPosition()
	// Keep player near center of viewport
	targetX := px - visW/2
	targetY := py - visH/2

	// Clamp viewport to map bounds (top-left origin)
	maxX := max(0, mapW-visW)
	maxY := max(0, mapH-visH)
	if targetX < 0 {
		targetX = 0
	}
	if targetY < 0 {
		targetY = 0
	}
	if targetX > maxX {
		targetX = maxX
	}
	if targetY > maxY {
		targetY = maxY
	}
	gr.viewX, gr.viewY = targetX, targetY
}

// renderMap draws the tile map into the grid, clipped within borders
func (gr *GameRenderer) renderMap(grid [][]rune) {
	if gr.tileMap == nil {
		return
	}
	visW := gr.width - 2
	visH := gr.height - 2
	mapW := gr.tileMap.Width
	mapH := gr.tileMap.Height

	if mapW <= 0 || mapH <= 0 {
		return
	}

	// Compute uniform integer scale to fit map when smaller than viewport
	sX := visW / mapW
	sY := visH / mapH
	gr.scale = 2
	if sX > 1 && sY > 1 {
		if sX < sY {
			gr.scale = sX
		} else {
			gr.scale = sY
		}
	}

	if gr.scale == 1 {
		// No scaling: draw the visible slice via viewport
		for y := 0; y < visH; y++ {
			mapY := gr.viewY + y
			for x := 0; x < visW; x++ {
				mapX := gr.viewX + x
				ch := gr.tileMap.At(mapX, mapY)
				if ch == 0 {
					ch = ' '
				}
				sy, sx := y+1, x+1
				if sy >= 1 && sy < gr.height-1 && sx >= 1 && sx < gr.width-1 {
					grid[sy][sx] = ch
				}
			}
		}
		return
	}

	// Scaling case: center the scaled map inside the viewport
	drawW := mapW * gr.scale
	drawH := mapH * gr.scale
	gr.offX = (visW - drawW) / 2
	gr.offY = (visH - drawH) / 2

	for my := 0; my < mapH; my++ {
		for mx := 0; mx < mapW; mx++ {
			ch := gr.tileMap.At(mx, my)
			if ch == 0 {
				ch = ' '
			}
			baseY := 1 + gr.offY + my*gr.scale
			baseX := 1 + gr.offX + mx*gr.scale
			for dy := 0; dy < gr.scale; dy++ {
				sy := baseY + dy
				if sy <= 0 || sy >= gr.height-1 {
					continue
				}
				for dx := 0; dx < gr.scale; dx++ {
					sx := baseX + dx
					if sx <= 0 || sx >= gr.width-1 {
						continue
					}
					grid[sy][sx] = ch
				}
			}
		}
	}
}

// renderPlayer draws the player sprite on the grid using viewport offset
func (gr *GameRenderer) renderPlayer(grid [][]rune, player *types.Player) {
	spriteLines := strings.Split(player.GetSprite(), "\n")
	playerX, playerY := player.GetPosition()
	for i, line := range spriteLines {
		var y, x int
		if gr.scale > 1 {
			y = 1 + gr.offY + (playerY-1)*gr.scale + i
		} else {
			y = (playerY - gr.viewY) + i + 1
		}
		if y >= 1 && y < gr.height-1 {
			for j, char := range line {
				if gr.scale > 1 {
					x = 1 + gr.offX + (playerX-1)*gr.scale + j
				} else {
					x = (playerX - gr.viewX) + j + 1
				}
				if x >= 1 && x < gr.width-1 {
					grid[y][x] = char
				}
			}
		}
	}
}

// gridToString converts the grid to a string efficiently
func (gr *GameRenderer) gridToString(grid [][]rune) string {
	var builder strings.Builder
	builder.Grow(gr.width * gr.height)
	for _, row := range grid {
		builder.WriteString(string(row))
		builder.WriteString("\n")
	}
	return strings.TrimRight(builder.String(), "\n")
}

// Extension points for future systems can be added here when needed.

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// (duplicate UpdateSize removed)
