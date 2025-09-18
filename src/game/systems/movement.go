package systems

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/game/types"
)

// MovementSystem handles player movement and collision detection
type MovementSystem struct {
	// Movement configuration
}

// NewMovementSystem creates a new movement system instance
func NewMovementSystem() *MovementSystem {
	return &MovementSystem{}
}

// spriteFootprintTiles returns the collision footprint in tiles.
// For now, keep it 1x1 to decouple from the rendered sprite.
func (ms *MovementSystem) spriteFootprintTiles(player *types.Player) (int, int) { return 4, 3 }

// isWalkable returns true if (x,y) is within map bounds and not a wall.
// x,y are 1-based player coordinates. When tm is nil, treat all tiles as walkable.
func (ms *MovementSystem) isWalkable(tm *types.TileMap, x, y int) bool {
	if x < 1 || y < 1 {
		return false
	}
	if tm == nil {
		return true
	}
	if x > tm.Width || y > tm.Height {
		return false
	}
	ch := tm.At(x-1, y-1)
	return !config.IsMapWall(ch)
}

// isWalkableRect checks a rectangle region of size (w x h) at top-left (x,y) for collisions and bounds.
// Returns true if all covered tiles are walkable.
func (ms *MovementSystem) isWalkableRect(tm *types.TileMap, x, y, w, h int) bool {
	if x < 1 || y < 1 || w < 1 || h < 1 {
		return false
	}
	if tm == nil {
		return true
	}
	if x+w-1 > tm.Width || y+h-1 > tm.Height {
		return false
	}
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			if !ms.isWalkable(tm, xx, yy) {
				return false
			}
		}
	}
	return true
}

// EnsureValidSpawn adjusts the player's position if it's inside a wall or out of bounds.
// It tries to find the nearest walkable tile using a simple expanding diamond search.
func (ms *MovementSystem) EnsureValidSpawn(player *types.Player, tm *types.TileMap) {
	if player == nil {
		return
	}
	px, py := player.Pos.X, player.Pos.Y
	wTiles, hTiles := ms.spriteFootprintTiles(player)
	if ms.isWalkableRect(tm, px, py, wTiles, hTiles) {
		return
	}

	// Fallback starting point
	startX, startY := 1, 1
	if tm != nil {
		// Clamp starting guess within map
		if px >= 1 && px <= tm.Width {
			startX = px
		}
		if py >= 1 && py <= tm.Height {
			startY = py
		}
	}

	// Expanding search radius for nearest walkable tile
	maxRadius := 50
	if tm != nil {
		if tm.Width+tm.Height > 0 {
			if tm.Width > tm.Height {
				maxRadius = tm.Width + 2
			} else {
				maxRadius = tm.Height + 2
			}
		}
	}
	for r := 0; r <= maxRadius; r++ {
		// Check center first when r==0
		if r == 0 {
			if ms.isWalkableRect(tm, startX, startY, wTiles, hTiles) {
				player.Pos.X, player.Pos.Y = startX, startY
				return
			}
			continue
		}
		// Diamond ring: (dx,dy) with |dx|+|dy| == r
		for dx := -r; dx <= r; dx++ {
			dy := r - abs(dx)
			candidates := [][2]int{
				{startX + dx, startY + dy},
				{startX + dx, startY - dy},
			}
			for _, c := range candidates {
				cx, cy := c[0], c[1]
				if ms.isWalkableRect(tm, cx, cy, wTiles, hTiles) {
					player.Pos.X, player.Pos.Y = cx, cy
					return
				}
			}
		}
	}
	// As a last resort, set to (1,1) if walkable; else leave as is
	if ms.isWalkableRect(tm, 1, 1, wTiles, hTiles) {
		player.Pos.X, player.Pos.Y = 1, 1
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// MovePlayer moves the player in the specified direction using map bounds and walls.
// The player is treated as occupying a single tile anchor; movement is blocked by wall tiles.
func (ms *MovementSystem) MovePlayer(player *types.Player, direction rune, tm *types.TileMap) bool {
	if player == nil {
		return false
	}

	oldX, oldY := player.Pos.X, player.Pos.Y
	wTiles, hTiles := ms.spriteFootprintTiles(player)

	dx, dy := 0, 0
	switch direction {
	case '↑':
		dy = -1
	case '↓':
		dy = 1
	case '←':
		dx = -1
	case '→':
		dx = 1
	default:
		return false
	}

	targetX := oldX + dx
	targetY := oldY + dy
	// Always keep coordinates at least 1-based (interior of the game view)
	if targetX < 1 {
		targetX = 1
	}
	if targetY < 1 {
		targetY = 1
	}

	if tm != nil {
		mapW, mapH := tm.Width, tm.Height
		// Clamp to map bounds accounting for sprite footprint
		maxX := mapW - wTiles + 1
		maxY := mapH - hTiles + 1
		if maxX < 1 {
			maxX = 1
		}
		if maxY < 1 {
			maxY = 1
		}
		if targetX > maxX {
			targetX = maxX
		}
		if targetY > maxY {
			targetY = maxY
		}

		// Check wall collision for the footprint rectangle
		if !ms.isWalkableRect(tm, targetX, targetY, wTiles, hTiles) {
			return false
		}
	}

	player.Pos.X = targetX
	player.Pos.Y = targetY
	return player.Pos.X != oldX || player.Pos.Y != oldY
}

// ValidatePosition checks if a position is within the game bounds
// ValidatePosition can be used for additional checks; here we simply ensure x,y are positive.
func (ms *MovementSystem) ValidatePosition(x, y, width, height int) bool {
	return x >= 1 && y >= 1
}

// GetPlayerPosition returns the current player position
func (ms *MovementSystem) GetPlayerPosition(player *types.Player) (int, int) {
	return player.GetPosition()
}
