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
	// inner viewport rectangle (borders will be drawn around this)
	innerX int // top-left X of viewport border in screen grid
	innerY int // top-left Y of viewport border in screen grid
	innerW int // viewport interior width (in cells)
	innerH int // viewport interior height (in cells)
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

	// Compute inner viewport layout and adjust camera
	gr.computeLayout(player)

	// Initialize draw state and grid
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
	if gr.innerW <= 0 || gr.innerH <= 0 {
		return
	}
	left := gr.innerX
	top := gr.innerY
	right := gr.innerX + gr.innerW + 1
	bottom := gr.innerY + gr.innerH + 1

	// Top and bottom
	for x := left; x <= right; x++ {
		if top >= 0 && top < gr.height {
			switch {
			case x == left:
				grid[top][x] = '┌'
			case x == right:
				grid[top][x] = '┐'
			default:
				grid[top][x] = '─'
			}
		}
		if bottom >= 0 && bottom < gr.height {
			switch {
			case x == left:
				grid[bottom][x] = '└'
			case x == right:
				grid[bottom][x] = '┘'
			default:
				grid[bottom][x] = '─'
			}
		}
	}
	// Sides
	for y := top + 1; y < bottom; y++ {
		if y >= 0 && y < gr.height {
			if left >= 0 && left < gr.width {
				grid[y][left] = '│'
			}
			if right >= 0 && right < gr.width {
				grid[y][right] = '│'
			}
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

// computeLayout decides whether to center a smaller viewport around the map
// or use the full available area with scrolling for large maps.
func (gr *GameRenderer) computeLayout(player *types.Player) {
	visW := max(1, gr.width-2)
	visH := max(1, gr.height-2)
	mapW, mapH := 0, 0
	if gr.tileMap != nil {
		mapW, mapH = gr.tileMap.Width, gr.tileMap.Height
	}

	// Small-map mode: center viewport to match map size and show entire map
	if mapW > 0 && mapH > 0 && mapW <= visW && mapH <= visH {
		gr.viewX, gr.viewY = 0, 0
		gr.innerW, gr.innerH = mapW, mapH
		totalW := gr.innerW + 2
		totalH := gr.innerH + 2
		gr.innerX = (gr.width - totalW) / 2
		gr.innerY = (gr.height - totalH) / 2
		if gr.innerX < 0 {
			gr.innerX = 0
		}
		if gr.innerY < 0 {
			gr.innerY = 0
		}
		return
	}

	// Large-map mode: full-screen interior viewport with scrolling
	gr.innerX, gr.innerY = 0, 0
	gr.innerW, gr.innerH = visW, visH
	gr.updateViewport(player)
}

// renderMap draws the tile map into the grid, clipped within borders
func (gr *GameRenderer) renderMap(grid [][]rune) {
	if gr.tileMap == nil {
		return
	}
	// Draw visible slice inside inner viewport (no scaling)
	for y := 0; y < gr.innerH; y++ {
		mapY := gr.viewY + y
		for x := 0; x < gr.innerW; x++ {
			mapX := gr.viewX + x
			ch := gr.tileMap.At(mapX, mapY)
			if ch == 0 {
				ch = ' '
			}
			sy := gr.innerY + 1 + y
			sx := gr.innerX + 1 + x
			if sy >= 0 && sy < gr.height && sx >= 0 && sx < gr.width {
				grid[sy][sx] = ch
			}
		}
	}
}

// renderPlayer draws the player sprite on the grid using viewport offset
func (gr *GameRenderer) renderPlayer(grid [][]rune, player *types.Player) {
	if player == nil {
		return
	}
	playerX, playerY := player.GetPosition()
	spriteLines := strings.Split(player.GetSprite(), "\n")
	for i, line := range spriteLines {
		y := gr.innerY + 1 + (playerY - gr.viewY - 1) + i
		if y >= 0 && y < gr.height {
			for j, char := range line {
				x := gr.innerX + 1 + (playerX - gr.viewX - 1) + j
				if x >= 0 && x < gr.width {
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
