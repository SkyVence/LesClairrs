// Package entities provides factory functions for creating and initializing
// game entities in ProjectRed RPG.
//
// This package handles the creation of complex game objects with proper
// initialization and setup:
//   - Player creation with class configuration
//   - World loading and initialization
//   - Entity lifecycle management
//
// Factory functions ensure consistent object creation and hide complex
// initialization logic from the rest of the application.
//
// Example usage:
//
//	player := entities.NewPlayer("Sam", class, position)
//	world := entities.NewWorld(1)
package entities

import "projectred-rpg.com/game/types"

// NewPlayer creates a new player with the specified name, class, and position
func NewPlayer(name string, class types.Class, pos types.Position) *types.Player {
	stats := types.PlayerStats{
		Level:        1,
		Exp:          0,
		NextLevelExp: 100,
		Force:        class.Force,
		Speed:        class.Speed,
		Defense:      class.Defense,
		Accuracy:     class.Accuracy,
		MaxHP:        class.MaxHP,
		CurrentHP:    class.MaxHP,
	}

	player := &types.Player{
		Name:      name,
		Class:     class,
		Stats:     stats,
		Pos:       pos,
		Inventory: make([]types.Item, 0, 10),
		Implants:  [5]types.Implant{},
		MaxInv:    10,
	}

	// Set the default sprite
	player.SetSprite(types.CreateStickManSprite())

	return player
}
