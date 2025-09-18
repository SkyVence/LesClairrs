package systems

import (
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

// CombatSystem handles combat calculations and interactions
type CombatSystem struct {
	CurrentCombatState  types.CombatState
	PreviousCombatState types.CombatState
	CurrentEnemy        *entities.Enemy
	locManager          *engine.LocalizationManager
	spawnerSystem       *SpawnerSystem
}

// NewCombatSystem creates a new combat system instance
func NewCombatSystem(initialState types.CombatState, locManager *engine.LocalizationManager, spawnerSystem *SpawnerSystem) *CombatSystem {
	return &CombatSystem{
		CurrentCombatState:  initialState,
		PreviousCombatState: initialState,
		locManager:          locManager,
		spawnerSystem:       spawnerSystem,
	}
}

func (cs *CombatSystem) ChangeCombatState(newState types.CombatState) {
	cs.PreviousCombatState = cs.CurrentCombatState
	cs.CurrentCombatState = newState
}

func (cs *CombatSystem) RevertCombatState() {
	cs.CurrentCombatState, cs.PreviousCombatState = cs.PreviousCombatState, cs.CurrentCombatState
}

func (cs *CombatSystem) IsInCombat() bool {
	return cs.CurrentCombatState != types.Idle && cs.CurrentCombatState != types.Dead
}

func (cs *CombatSystem) EnterCombat(e *entities.Enemy, p *types.Player) {
	cs.CurrentEnemy = e
	cs.ChangeCombatState(types.PlayerTurn)
}

func (cs *CombatSystem) AiAttack(e entities.Enemy, p *types.Player) {
	cs.ChangeCombatState(types.PlayerTurn)
}

func (cs *CombatSystem) PlayerAttack(e *entities.Enemy, p *types.Player) {
	damage := p.CalculateDamage(e.Defense)
	defeated := e.TakeDamage(damage)
	if defeated {
		cs.ChangeCombatState(types.Victory)
		p.AddExperience(e.ExpReward)
	} else {
		cs.ChangeCombatState(types.EnemyTurn)
	}
}

func (cs *CombatSystem) UseConsumable(item types.Item, p *types.Player) string {
	// Placeholder for using consumable items
	return cs.locManager.Text(item.Name) + " used." // Example return message
}

// IsEnemyDefeated for checking if enemy is defeated
func (cs *CombatSystem) IsEnemyDefeated(e entities.Enemy, p *types.Player) bool {
	return e.CurrentHP <= 0
}

// IsPlayerDefeated checks if the player is defeated
func (cs *CombatSystem) IsPlayerDefeated(player *types.Player) bool {
	return player.Stats.CurrentHP <= 0
}

// CheckForCombatEngagement checks if the player is close enough to any enemy to start combat
// Returns the enemy to engage with, or nil if no engagement
func (cs *CombatSystem) CheckForCombatEngagement(player *types.Player) *entities.Enemy {
	if cs.IsInCombat() {
		return nil // Already in combat
	}

	const combatRange = 3.0 // 2 character distance as specified
	return cs.spawnerSystem.CheckPlayerProximity(player.Pos, combatRange)
}

// TryEngageCombat attempts to engage combat if player is within range
// Returns true if combat was initiated
func (cs *CombatSystem) TryEngageCombat(player *types.Player) bool {
	enemy := cs.CheckForCombatEngagement(player)
	if enemy != nil {
		cs.EnterCombat(enemy, player)
		return true
	}
	return false
}

// GetCurrentEnemy returns the current enemy being fought
func (cs *CombatSystem) GetCurrentEnemy() *entities.Enemy {
	return cs.CurrentEnemy
}

// ExitCombat ends the current combat encounter
func (cs *CombatSystem) ExitCombat() {
	cs.CurrentEnemy = nil
	cs.ChangeCombatState(types.Idle)
}
