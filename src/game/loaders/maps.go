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
	return types.NewTileMap(lines)
}
