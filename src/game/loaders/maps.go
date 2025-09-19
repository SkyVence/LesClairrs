package loaders

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"projectred-rpg.com/game/types"
)

// LoadStageMap tries to load a map file for a given world and stage.
// Files are expected under assets/levels as world-<id>_stage-<nb>.map.
// Returns nil if not found or on error (caller can fallback to empty background).
func LoadStageMap(worldID, stageNb int) *types.TileMap {
	// Construct filename like: assets/levels/world-1_stage-1.map
	fileName := fmt.Sprintf("world-%d_stage-%d.map", worldID, stageNb)
	mapPath := filepath.Join("assets", "levels", fileName)

	f, err := os.Open(mapPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0, 64)
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r"))
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	if len(lines) == 0 {
		return nil
	}

	tm := types.NewTileMap(lines)

	// Apply custom transition zone coordinates if available
	if tm != nil {
		if x, y, w, h, hasCustom := getCustomTransitionZone(worldID, stageNb); hasCustom {
			tm.SetCustomTransitionZone(x, y, w, h)
		}
	}

	return tm
}

// getCustomTransitionZone returns custom transition zone coordinates for specific world/stage combinations
func getCustomTransitionZone(worldID, stageNb int) (x, y, width, height int, hasCustom bool) {
	switch worldID {
	case 1:
		switch stageNb {
		case 1:
			return 4, 1, 16, 2, true // Custom position for world 1, stage 1
		case 2:
			return 5, 15, 3, 2, true // Custom position for world 1, stage 2
		case 3:
			return 20, 5, 2, 3, true // Custom position for world 1, stage 3
		}
	case 2:
		switch stageNb {
		case 1:
			return 15, 12, 2, 2, true // Custom position for world 2, stage 1
		case 2:
			return 8, 6, 3, 2, true // Custom position for world 2, stage 2
		}
	}
	// Return false if no custom coordinates defined for this world/stage
	return 0, 0, 0, 0, false
}
