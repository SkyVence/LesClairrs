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

import "time"

// Game balance constants
const (
	DefaultPlayerHealth = 100
	MaxInventorySize    = 10
	BaseExpRequirement  = 100
	ExpGrowthRate       = 1.2

	// Animation constants
	DefaultAnimationSpeed = 200 * time.Millisecond
	TickDuration          = 16 * time.Millisecond // ~60 FPS

	// Display constants
	DefaultTerminalWidth  = 80
	DefaultTerminalHeight = 24
	HUDHeight             = 5
)

// ClassConfig represents a character class configuration
type ClassConfig struct {
	Name        string
	Description string
	MaxHP       int
	Force       int
	Speed       int
	Defense     int
	Accuracy    int
}

// Default classes available in the game
var DefaultClasses = map[string]ClassConfig{
	"D0C": {
		Name:        "D0C",
		Description: "Un robot intelligent, précis et polyvalent.",
		MaxHP:       90,
		Force:       10,
		Speed:       12,
		Defense:     10,
		Accuracy:    22,
	},
	"APP": {
		Name:        "APP",
		Description: "Un robot furtif, rapide et précis.",
		MaxHP:       80,
		Force:       14,
		Speed:       22,
		Defense:     8,
		Accuracy:    18,
	},
	"CYBER_SAMURAI": {
		Name:        "Cyber-Samurai",
		Description: "A swift and deadly warrior, excelling in close combat and agility.",
		MaxHP:       100,
		Force:       15,
		Speed:       12,
		Defense:     8,
		Accuracy:    15,
	},
}
