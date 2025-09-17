package entities

import (
	"projectred-rpg.com/game/loaders"
	"projectred-rpg.com/game/types"
)

func NewWorld(worldID int) *types.World {
	if err := loaders.LoadWorlds(); err != nil {
		return &types.World{WorldID: worldID}
	}

	if world, exists := loaders.GetWorld(worldID); exists {
		return &world
	}
	return &types.World{WorldID: worldID}
}
