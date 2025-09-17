// Package game provides the main game logic coordination for ProjectRed RPG.
//
// This package serves as the central coordinator that orchestrates all game systems,
// manages game state, and handles the main game loop. It integrates:
//   - Entity management (players, worlds, stages)
//   - System coordination (combat, inventory, movement)
//   - Game flow and progression
//   - State transitions and location tracking
//
// The Game struct contains all necessary systems and provides methods for
// game progression, location management, and state queries.
//
// Example usage:
//
//	game := NewGameInstance(selectedClass)
//	location, worldID := game.CurrentLocation()
//	success := game.Advance()
package game

import (
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/loaders"
	"projectred-rpg.com/game/systems"
	"projectred-rpg.com/game/types"
)

// Game represents the main game state and coordinates all game systems.
// It serves as the central hub that manages:
//   - Current player state and progression
//   - World and stage navigation
//   - Integration of all game systems (combat, inventory, movement)
//   - Game flow and state transitions
type Game struct {
	Player       *types.Player // Current player instance with stats and inventory
	CurrentWorld *types.World  // Currently active world containing stages
	CurrentStage *types.Stage  // Current stage within the active world

	// Game systems - modular components handling specific game logic
	Combat    *systems.CombatSystem    // Handles damage calculations and battle mechanics
	Inventory *systems.InventorySystem // Manages item operations and equipment
	Movement  *systems.MovementSystem  // Processes player movement and collision detection
}

// NewGameInstance creates a new game with the specified character class.
// This function initializes all game systems and creates the starting game state:
//   - Creates a player with the selected class and default stats
//   - Loads the first world and sets the starting stage
//   - Initializes all game systems (combat, inventory, movement)
//
// Parameters:
//
//	selectedClass: The character class that determines base stats and abilities
//
// Returns:
//
//	*Game: Fully initialized game instance ready for play
//
// Example:
//
//	class := config.DefaultClasses["CYBER_SAMURAI"]
//	game := NewGameInstance(class)
func NewGameInstance(selectedClass types.Class) *Game {
	world := NewWorld(1)

	// Hardcoded for now --> Probably will be changed if implementing save/load system
	player := entities.NewPlayer("Sam", selectedClass, types.Position{X: 1, Y: 1})
	return &Game{
		Player:       player,
		CurrentWorld: world,
		CurrentStage: &world.Stages[0],
		// Initialize systems
		Combat:    systems.NewCombatSystem(),
		Inventory: systems.NewInventorySystem(),
		Movement:  systems.NewMovementSystem(),
	}
}

// NewWorld loads or creates a world by its ID.
// This function attempts to load world data from the cache, falling back to
// creating an empty world if loading fails.
//
// Parameters:
//
//	WorldID: The unique identifier for the world to load
//
// Returns:
//
//	*types.World: The loaded world data, or an empty world with the specified ID
//
// The function ensures worlds are loaded before attempting retrieval and
// provides graceful fallback behavior for missing or corrupted world data.
func NewWorld(WorldID int) *types.World {
	// Ensure worlds are loaded
	if err := loaders.LoadWorlds(); err != nil {
		return &types.World{WorldID: WorldID}
	}

	if world, exists := loaders.GetWorld(WorldID); exists {
		return &world
	}
	return &types.World{WorldID: WorldID}
}

// CurrentLocation returns a human-readable location string for the HUD.
// This method provides the current location information for display purposes,
// preferring stage names when available and falling back to world names.
//
// Returns:
//
//	string: Human-readable location name (stage name or world name)
//	int: World ID for additional context (-1 if game state is invalid)
//
// The returned information is suitable for display in the game's HUD or
// status indicators to show the player's current position in the game world.
func (g *Game) CurrentLocation() (string, int) {
	if g == nil || g.CurrentWorld == nil {
		return "", -1
	}
	if g.CurrentStage != nil && g.CurrentStage.Name != "" {
		// Prefer stage name when available
		return g.CurrentStage.Name, g.CurrentWorld.WorldID
	}
	return g.CurrentWorld.Name, g.CurrentWorld.WorldID
}

// PeekNext returns information about the next stage or world without mutating state.
// It returns a user-facing name for the upcoming location and the target world ID,
// plus a boolean indicating whether a next location exists.
func (g *Game) PeekNext() (string, int, bool) {
	if g == nil || g.CurrentWorld == nil || g.CurrentStage == nil {
		return "", -1, false
	}
	// Next stage in current world?
	if next := g.CurrentStage.AdvanceToNextStage(g.CurrentWorld); next != nil {
		name := next.Name
		if name == "" {
			name = g.CurrentWorld.Name
		}
		return name, g.CurrentWorld.WorldID, true
	}
	// Otherwise next world, first stage
	if nw := g.CurrentWorld.NextWorld(); nw != nil {
		if len(nw.Stages) > 0 {
			name := nw.Stages[0].Name
			if name == "" {
				name = nw.Name
			}
			return name, nw.WorldID, true
		}
	}
	return "", -1, false
}

// Advance moves the game to the next stage; if no more stages exist,
// it loads the next world and moves to its first stage. Returns true on success.
func (g *Game) Advance() bool {
	if g == nil || g.CurrentWorld == nil || g.CurrentStage == nil {
		return false
	}
	if next := g.CurrentStage.AdvanceToNextStage(g.CurrentWorld); next != nil {
		g.CurrentStage = next
		return true
	}
	// Try next world
	if nw := g.CurrentWorld.NextWorld(); nw != nil {
		g.CurrentWorld = nw
		if len(nw.Stages) > 0 {
			g.CurrentStage = &nw.Stages[0]
			return true
		}
	}
	return false
}
