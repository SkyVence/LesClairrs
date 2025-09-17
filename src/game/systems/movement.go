package systems

import (
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

// MovePlayer moves the player in the specified direction within bounds
func (ms *MovementSystem) MovePlayer(player *types.Player, direction rune, width, height int) bool {
	oldX, oldY := player.Pos.X, player.Pos.Y

	player.Move(direction, width, height)

	// Return true if position actually changed
	return player.Pos.X != oldX || player.Pos.Y != oldY
}

// ValidatePosition checks if a position is within the game bounds
func (ms *MovementSystem) ValidatePosition(x, y, width, height int) bool {
	return x >= 1 && x < width-4 && y >= 1 && y < height-4
}

// GetPlayerPosition returns the current player position
func (ms *MovementSystem) GetPlayerPosition(player *types.Player) (int, int) {
	return player.GetPosition()
}
