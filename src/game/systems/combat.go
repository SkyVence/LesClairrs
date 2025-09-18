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
	locManager          *engine.LocalizationManager
}

// NewCombatSystem creates a new combat system instance
func NewCombatSystem(initialState types.CombatState, locManager *engine.LocalizationManager) *CombatSystem {
	return &CombatSystem{
		CurrentCombatState:  initialState,
		PreviousCombatState: initialState,
		locManager:          locManager,
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

func (cs *CombatSystem) EnterCombat(e entities.Enemy, p *types.Player) {
	cs.ChangeCombatState(types.Idle)
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
