package systems

import (
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

type SpawnerSystem struct {
	ActiveEnemies []*entities.Enemy
	Stage         *types.Stage
}

// NewSpawnerSystem creates a new spawner system
func NewSpawnerSystem() *SpawnerSystem {
	return &SpawnerSystem{
		ActiveEnemies: make([]*entities.Enemy, 0),
	}
}

// LoadStage loads enemies from a stage definition
func (ss *SpawnerSystem) LoadStage(stage *types.Stage) {
	ss.Stage = stage
	ss.ActiveEnemies = make([]*entities.Enemy, 0)

	// Create enemies from the stage's enemy spawn data
	for _, enemySpawn := range stage.Enemies {
		enemy := entities.NewEnemy(entities.Enemy{
			Name:      enemySpawn.Name,
			Force:     enemySpawn.Force,
			Speed:     enemySpawn.Speed,
			Defense:   enemySpawn.Defense,
			Accuracy:  enemySpawn.Accuracy,
			MaxHP:     enemySpawn.MaxHP,
			CurrentHP: enemySpawn.CurrentHP,
			ExpReward: enemySpawn.ExpReward,
			Sprite:    enemySpawn.Sprite,
			Position:  enemySpawn.Position,
		})
		ss.ActiveEnemies = append(ss.ActiveEnemies, enemy)
	}
}

// GetActiveEnemies returns all active (alive) enemies
func (ss *SpawnerSystem) GetActiveEnemies() []*entities.Enemy {
	activeEnemies := make([]*entities.Enemy, 0)
	for _, enemy := range ss.ActiveEnemies {
		if enemy.IsAlive {
			activeEnemies = append(activeEnemies, enemy)
		}
	}
	return activeEnemies
}

// RemoveDefeatedEnemies removes all defeated enemies from the active list
func (ss *SpawnerSystem) RemoveDefeatedEnemies() {
	activeEnemies := make([]*entities.Enemy, 0)
	for _, enemy := range ss.ActiveEnemies {
		if enemy.IsAlive {
			activeEnemies = append(activeEnemies, enemy)
		}
	}
	ss.ActiveEnemies = activeEnemies
}

// CheckPlayerProximity checks if the player is within combat range of any enemy
func (ss *SpawnerSystem) CheckPlayerProximity(playerPos types.Position, combatRange float64) *entities.Enemy {
	for _, enemy := range ss.GetActiveEnemies() {
		if enemy.IsWithinRange(playerPos, combatRange) {
			return enemy
		}
	}
	return nil
}

// GetEnemyCount returns the total number of active enemies
func (ss *SpawnerSystem) GetEnemyCount() int {
	return len(ss.GetActiveEnemies())
}

// IsStageCleared returns true if all enemies are defeated
func (ss *SpawnerSystem) IsStageCleared() bool {
	return ss.GetEnemyCount() == 0
}
