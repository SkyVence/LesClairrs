// Package loaders handles asset loading and caching for ProjectRed RPG.
//
// This package provides functionality for loading game data from various sources:
//   - JSON file parsing for worlds and levels
//   - Asset caching for performance
//   - Data validation and error handling
//   - Thread-safe access to cached data
//
// The package uses a repository pattern with caching to provide fast access
// to game data while abstracting the underlying storage format.
//
// Example usage:
//
//	err := loaders.LoadWorlds()
//	world, exists := loaders.GetWorld(1)
package loaders

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"projectred-rpg.com/game/types"
)

var (
	worldCache   map[int]types.World
	worldMutex   sync.RWMutex
	worldsLoaded bool = false
)

// LoadWorlds loads all worlds from JSON files in assets/levels directory
func LoadWorlds() error {
	worldMutex.Lock()
	defer worldMutex.Unlock()

	if worldsLoaded {
		return nil // Already loaded
	}

	worldCache = make(map[int]types.World)

	// Get the path to the assets/levels directory
	levelsPath := filepath.Join("assets", "levels")

	// Read all JSON files in the directory
	err := filepath.WalkDir(levelsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		// Read the JSON file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read world file %s: %w", path, err)
		}

		// Parse the JSON into a World struct
		var world types.World
		if err := json.Unmarshal(data, &world); err != nil {
			return fmt.Errorf("failed to parse world file %s: %w", path, err)
		}

		// Set WorldID for all stages in this world
		for i := range world.Stages {
			world.Stages[i].WorldID = world.WorldID
		}

		// Store the world in the cache using its WorldID
		worldCache[world.WorldID] = world

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load worlds: %w", err)
	}

	worldsLoaded = true
	return nil
}

// GetWorld retrieves a world by its ID
func GetWorld(worldID int) (types.World, bool) {
	worldMutex.RLock()
	defer worldMutex.RUnlock()

	world, exists := worldCache[worldID]
	return world, exists
}

// GetAllWorlds returns a copy of all loaded worlds
func GetAllWorlds() map[int]types.World {
	worldMutex.RLock()
	defer worldMutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[int]types.World)
	for k, v := range worldCache {
		result[k] = v
	}
	return result
}

// GetStage retrieves a specific stage from a world
func GetStage(worldID, stageNb int) (types.Stage, bool) {
	worldMutex.RLock()
	defer worldMutex.RUnlock()

	world, exists := worldCache[worldID]
	if !exists {
		return types.Stage{}, false
	}

	for _, stage := range world.Stages {
		if stage.StageNb == stageNb {
			return stage, true
		}
	}
	return types.Stage{}, false
}
