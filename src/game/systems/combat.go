// Package systems implements modular game logic systems for ProjectRed RPG.
//
// This package follows a systems-based architecture where each system handles
// a specific aspect of game logic:
//   - CombatSystem: Damage calculations, battle mechanics
//   - InventorySystem: Item management, equipment handling
//   - MovementSystem: Player movement, collision detection
//
// Systems operate on types from the game/types package and provide
// clean separation of concerns. Each system is independently testable
// and can be easily extended or replaced.
//
// Example usage:
//
//	combat := systems.NewCombatSystem()
//	damage := combat.PlayerAttacksEnemy(player, enemy)
//	defeated := combat.IsEnemyDefeated(enemy)
package systems

import (
	"projectred-rpg.com/game/types"
)

// CombatSystem handles combat calculations and interactions
type CombatSystem struct {
	// Combat state and configuration
}

// NewCombatSystem creates a new combat system instance
func NewCombatSystem() *CombatSystem {
	return &CombatSystem{}
}

// CalculateDamage calculates damage between an attacker and target
func (cs *CombatSystem) CalculateDamage(attackerForce, targetDefense int) int {
	damage := attackerForce - targetDefense
	if damage < 0 {
		damage = 0
	}
	return damage
}

// PlayerAttacksEnemy handles when a player attacks an enemy
func (cs *CombatSystem) PlayerAttacksEnemy(player *types.Player, enemy *types.Enemy) int {
	damage := cs.CalculateDamage(player.Stats.Force, enemy.Defense)
	enemy.CurrentHP -= damage
	if enemy.CurrentHP < 0 {
		enemy.CurrentHP = 0
	}
	return damage
}

// EnemyAttacksPlayer handles when an enemy attacks a player
func (cs *CombatSystem) EnemyAttacksPlayer(enemy *types.Enemy, player *types.Player) int {
	damage := cs.CalculateDamage(enemy.Force, player.Stats.Defense)
	player.Stats.CurrentHP -= damage
	if player.Stats.CurrentHP < 0 {
		player.Stats.CurrentHP = 0
	}
	return damage
}

// IsEnemyDefeated checks if an enemy is defeated
func (cs *CombatSystem) IsEnemyDefeated(enemy *types.Enemy) bool {
	return enemy.CurrentHP <= 0
}

// IsPlayerDefeated checks if the player is defeated
func (cs *CombatSystem) IsPlayerDefeated(player *types.Player) bool {
	return player.Stats.CurrentHP <= 0
}
