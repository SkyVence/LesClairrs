// Package config provides centralized configuration management for ProjectRed RPG.
//
// This package contains all game configuration including:
//   - Game balance constants (health, inventory, experience)
//   - Animation and timing settings
//   - Display configuration
//   - Character class definitions
//
// The configuration is designed to be easily modifiable for game tuning
// and supports different game modes or difficulty levels.
//
// Example usage:
//
//	health := config.DefaultPlayerHealth
//	class := config.DefaultClasses["CYBER_SAMURAI"]
package config

import (
	"projectred-rpg.com/game/types"
)

// Default classes available in the game
func GetDefaultClasses() []types.Class {
	return []types.Class{
		{
			Name:        "ui.class.doc.name",
			Description: "ui.class.doc.description",
			MaxHP:       90,
			Force:       10,
			Speed:       12,
			Defense:     10,
			Accuracy:    22,
		},
		{
			Name:        "ui.class.app.name",
			Description: "ui.class.app.description",
			MaxHP:       80,
			Force:       14,
			Speed:       22,
			Defense:     8,
			Accuracy:    18,
		},
		{
			Name:        "ui.class.per.name",
			Description: "ui.class.per.description",
			MaxHP:       100,
			Force:       15,
			Speed:       12,
			Defense:     8,
			Accuracy:    15,
		},
	}
}
