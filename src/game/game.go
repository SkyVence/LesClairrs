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
//	gameRender := GameModel()
//	location, worldID := gameRender.CurrentLocation()
package game

import (
	"fmt"

	"projectred-rpg.com/engine"
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
	//Combat    *systems.CombatSystem    // Handles damage calculations and battle mechanics
	Inventory  *systems.InventorySystem  // Manages item operations and equipment
	Movement   *systems.MovementSystem   // Processes player movement and collision detection
	LevelIntro *systems.LevelIntroSystem // Handles level introduction dialogues

	// Game state
	language     string
	pendingStage *types.Stage // Stage to load after intro completes
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
func NewGameInstance(selectedClass types.Class, language string) *Game {
	world := NewWorld(1)
	player := entities.NewPlayer("Sam", selectedClass, types.Position{X: 1, Y: 1})

	// Create level intro system
	levelIntro := systems.NewLevelIntroSystem(language)
	if err := levelIntro.LoadLocalization(); err != nil {
		panic(fmt.Sprintf("Failed to load localization for language '%s': %v", language, err))
	}

	return &Game{
		Player:       player,
		CurrentWorld: world,
		CurrentStage: &world.Stages[0],
		Inventory:    systems.NewInventorySystem(),
		Movement:     systems.NewMovementSystem(),
		LevelIntro:   levelIntro, // AJOUTEZ CETTE LIGNE
		language:     language,
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

// LoadStage loads a stage with optional introduction
func (g *Game) LoadStage(worldID, stageID int) {

	// Create filename in the format your system expects
	filename := fmt.Sprintf("world-%d_stage-%d.map", worldID, stageID)

	// Find the target stage
	var targetStage *types.Stage
	if g.CurrentWorld != nil && g.CurrentWorld.WorldID == worldID {
		// Same world, find stage
		for i, stage := range g.CurrentWorld.Stages {
			if stage.StageNb == stageID {
				targetStage = &g.CurrentWorld.Stages[i]
				break
			}
		}
	} else {
		// Different world, load it first
		g.CurrentWorld = NewWorld(worldID)
		if len(g.CurrentWorld.Stages) > stageID-1 {
			targetStage = &g.CurrentWorld.Stages[stageID-1]
		}
	}

	if targetStage == nil {
		return // Stage not found
	}

	// Try to show intro
	if g.LevelIntro.ShowIntro(filename, 80, 24, func() {
		// Callback when intro is complete - actually load the stage
		g.actuallyLoadStage(targetStage)
	}) {
		// Intro found, store the stage to load after intro
		g.pendingStage = targetStage
	} else {
		// No intro, load stage directly
		g.actuallyLoadStage(targetStage)
	}
}

// actuallyLoadStage performs the actual stage loading
func (g *Game) actuallyLoadStage(stage *types.Stage) {
	g.CurrentStage = stage
	g.pendingStage = nil
	// Add any other stage loading logic here
}

// IsShowingIntro returns whether an intro is currently being shown
func (g *Game) IsShowingIntro() bool {
	return g.LevelIntro != nil && g.LevelIntro.IsActive()
}

// GameRender methods for accessing game state through the render interface

// CurrentLocation returns the current location information via GameRender.
func (gr *GameRender) CurrentLocation() (string, int) {
	if gr.gameInstance == nil {
		return "", -1
	}
	return gr.gameInstance.CurrentLocation()
}

// GetGameState returns the current game state.
func (gr *GameRender) GetGameState() systems.GameState {
	if gr.gameState == nil {
		return systems.GameState{}
	}
	return *gr.gameState
}

// GetPlayer returns the current player instance.
func (gr *GameRender) GetPlayer() *types.Player {
	if gr.gameInstance == nil {
		return nil
	}
	return gr.gameInstance.Player
}

// GetCurrentWorld returns the current world.
func (gr *GameRender) GetCurrentWorld() *types.World {
	if gr.gameInstance == nil {
		return nil
	}
	return gr.gameInstance.CurrentWorld
}

// GetCurrentStage returns the current stage.
func (gr *GameRender) GetCurrentStage() *types.Stage {
	if gr.gameInstance == nil {
		return nil
	}
	return gr.gameInstance.CurrentStage
}

// Main entry point function that creates and returns the GameRender model
// This replaces any previous main initialization and should be called by the engine
func InitGame() engine.Model {
	return GameModel()
}
