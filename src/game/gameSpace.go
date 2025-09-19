package game

import (
	"strings"

	"projectred-rpg.com/config"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

// Game space renderer with viewport/camera
type GameRenderer struct {
	width   int
	height  int
	tileMap *types.TileMap
	enemies []*entities.Enemy
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
	gr.renderEnemies(grid)
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

func (gr *GameRenderer) SetEnemies(enemies []*entities.Enemy) {
	gr.enemies = enemies
}

// ForceRefreshEnemies immediately updates the enemy list to reflect current spawner state
func (gr *GameRenderer) ForceRefreshEnemies(enemies []*entities.Enemy) {
	gr.enemies = make([]*entities.Enemy, len(enemies))
	copy(gr.enemies, enemies)
}

// RemoveEnemy removes a specific enemy from the renderer by reference
func (gr *GameRenderer) RemoveEnemy(enemyToRemove *entities.Enemy) bool {
	if gr.enemies == nil {
		return false
	}

	for i, enemy := range gr.enemies {
		if enemy == enemyToRemove {
			// Remove enemy by slicing
			gr.enemies = append(gr.enemies[:i], gr.enemies[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveEnemyByPosition removes an enemy at a specific position
func (gr *GameRenderer) RemoveEnemyByPosition(pos types.Position) bool {
	if gr.enemies == nil {
		return false
	}

	for i, enemy := range gr.enemies {
		enemyPos := enemy.GetPosition()
		if enemyPos.X == pos.X && enemyPos.Y == pos.Y {
			gr.enemies = append(gr.enemies[:i], gr.enemies[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveDeadEnemies removes all enemies that are not alive from the renderer
func (gr *GameRenderer) RemoveDeadEnemies() int {
	if gr.enemies == nil {
		return 0
	}

	originalCount := len(gr.enemies)
	aliveEnemies := make([]*entities.Enemy, 0, len(gr.enemies))

	for _, enemy := range gr.enemies {
		if enemy.IsAlive {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	gr.enemies = aliveEnemies
	return originalCount - len(gr.enemies)
}

// ClearAllEnemies removes all enemies from the renderer
func (gr *GameRenderer) ClearAllEnemies() {
	gr.enemies = nil
}

// GetEnemyCount returns the current number of enemies in the renderer
func (gr *GameRenderer) GetEnemyCount() int {
	if gr.enemies == nil {
		return 0
	}
	return len(gr.enemies)
}

// GetAliveEnemyCount returns the number of alive enemies
func (gr *GameRenderer) GetAliveEnemyCount() int {
	if gr.enemies == nil {
		return 0
	}

	count := 0
	for _, enemy := range gr.enemies {
		if enemy.IsAlive {
			count++
		}
	}
	return count
}

// GetEnemiesInArea returns all enemies within a rectangular area
func (gr *GameRenderer) GetEnemiesInArea(topLeft, bottomRight types.Position) []*entities.Enemy {
	if gr.enemies == nil {
		return nil
	}

	enemiesInArea := make([]*entities.Enemy, 0)
	for _, enemy := range gr.enemies {
		pos := enemy.GetPosition()
		if pos.X >= topLeft.X && pos.X <= bottomRight.X &&
			pos.Y >= topLeft.Y && pos.Y <= bottomRight.Y {
			enemiesInArea = append(enemiesInArea, enemy)
		}
	}
	return enemiesInArea
}

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
// Skips rendering outer wall characters (first/last row/column) but keeps them in map data
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

			// Check if this position is in an active transition zone
			if gr.tileMap.TransitionZone != nil && gr.tileMap.TransitionZone.Active &&
				gr.tileMap.TransitionZone.IsInZone(mapX, mapY) {
				ch = '◊' // Special character for transition zone
			}

			// Skip rendering outer walls (first/last row/column of the map)
			// but keep them in the map data for collision detection
			if gr.isOuterWall(mapX, mapY) {
				ch = ' ' // Render as empty space instead
			}

			sy := gr.innerY + 1 + y
			sx := gr.innerX + 1 + x
			if sy >= 0 && sy < gr.height && sx >= 0 && sx < gr.width {
				grid[sy][sx] = ch
			}
		}
	}
}

func (gr *GameRenderer) renderEnemies(grid [][]rune) {
	if gr.enemies == nil {
		return
	}

	// Default enemy sprite fallback
	defaultEnemySprite := ` ●  
/|\/
/ \`

	for _, enemy := range gr.enemies {
		if !enemy.IsAlive {
			continue
		}

		pos := enemy.GetPosition()
		enemyX, enemyY := pos.X, pos.Y

		// Check if enemy is visible in viewport
		if enemyX >= gr.viewX && enemyX < gr.viewX+gr.innerW &&
			enemyY >= gr.viewY && enemyY < gr.viewY+gr.innerH {

			// Use enemy's own sprite if available, otherwise use default
			enemySprite := enemy.Sprite
			if enemySprite == "" {
				enemySprite = defaultEnemySprite
			}

			// Render enemy sprite similar to player rendering
			spriteLines := strings.Split(enemySprite, "\n")
			for i, line := range spriteLines {
				y := gr.innerY + 1 + (enemyY - gr.viewY - 1) + i
				if y >= 0 && y < gr.height {
					for j, char := range line {
						x := gr.innerX + 1 + (enemyX - gr.viewX - 1) + j
						if x >= 0 && x < gr.width && char != ' ' {
							grid[y][x] = char
						}
					}
				}
			}
		}
	}
}

// isOuterWall checks if the given map coordinates are on the outer border
func (gr *GameRenderer) isOuterWall(mapX, mapY int) bool {
	if gr.tileMap == nil {
		return false
	}
	// Use the config function to check for outer walls
	return config.IsOuterWall(mapX, mapY, gr.tileMap.Width, gr.tileMap.Height)
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

// AddEnemy adds a new enemy to the renderer
func (gr *GameRenderer) AddEnemy(enemy *entities.Enemy) {
	if gr.enemies == nil {
		gr.enemies = make([]*entities.Enemy, 0)
	}
	gr.enemies = append(gr.enemies, enemy)
}

// FindEnemyByPosition finds an enemy at a specific position
func (gr *GameRenderer) FindEnemyByPosition(pos types.Position) *entities.Enemy {
	if gr.enemies == nil {
		return nil
	}

	for _, enemy := range gr.enemies {
		enemyPos := enemy.GetPosition()
		if enemyPos.X == pos.X && enemyPos.Y == pos.Y && enemy.IsAlive {
			return enemy
		}
	}
	return nil
}

// GetVisibleEnemies returns all enemies currently visible in the viewport
func (gr *GameRenderer) GetVisibleEnemies() []*entities.Enemy {
	if gr.enemies == nil {
		return nil
	}

	visibleEnemies := make([]*entities.Enemy, 0)
	for _, enemy := range gr.enemies {
		if !enemy.IsAlive {
			continue
		}

		pos := enemy.GetPosition()
		if pos.X >= gr.viewX && pos.X < gr.viewX+gr.innerW &&
			pos.Y >= gr.viewY && pos.Y < gr.viewY+gr.innerH {
			visibleEnemies = append(visibleEnemies, enemy)
		}
	}
	return visibleEnemies
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
