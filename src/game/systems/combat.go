package systems

import (
	"fmt"
	"math/rand"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
	"projectred-rpg.com/ui"
)

// CombatSystem handles combat calculations and interactions
type CombatSystem struct {
	CurrentCombatState  types.CombatState
	PreviousCombatState types.CombatState
	CurrentEnemy        *entities.Enemy
	locManager          *engine.LocalizationManager
	spawnerSystem       *SpawnerSystem
	combatUI           *ui.CombatUI
}

// NewCombatSystem creates a new combat system instance
func NewCombatSystem(initialState types.CombatState, locManager *engine.LocalizationManager, spawnerSystem *SpawnerSystem) *CombatSystem {
	return &CombatSystem{
		CurrentCombatState:  initialState,
		PreviousCombatState: initialState,
		locManager:          locManager,
		spawnerSystem:       spawnerSystem,
		combatUI:           nil, // Will be initialized later when renderer is available
	}
}

// SetRenderer sets the renderer and initializes the combat UI
func (cs *CombatSystem) SetRenderer(renderer engine.Renderer) {
	cs.combatUI = ui.NewCombatUI(renderer, cs.locManager)
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
	
	// Set up the combat UI if available
	if cs.combatUI != nil {
		cs.combatUI.SetCombatants(p, e)
		cs.combatUI.SetTurn(types.PlayerTurn)
		cs.combatUI.AddAction("System", "Combat", "", 0, fmt.Sprintf("Combat started against %s!", e.Name))
	}
}

func (cs *CombatSystem) AiAttack(e entities.Enemy, p *types.Player) {
	// Calculate AI damage
	baseDamage := e.Force
	damage := baseDamage - p.Stats.Defense
	if damage < 1 {
		damage = 1
	}
	
	// Add some randomness to AI attacks
	damage += rand.Intn(3) - 1 // -1 to +1 variation
	if damage < 1 {
		damage = 1
	}
	
	// Apply damage
	p.Stats.CurrentHP -= damage
	if p.Stats.CurrentHP < 0 {
		p.Stats.CurrentHP = 0
	}
	
	// Log the action
	message := fmt.Sprintf("%s attacks %s for %d damage!", e.Name, p.Name, damage)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(e.Name, "Attack", p.Name, damage, message)
	}
	
	// Check if player is defeated
	if cs.IsPlayerDefeated(p) {
		cs.ChangeCombatState(types.Dead)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.Dead)
			cs.combatUI.AddAction("System", "Result", "", 0, "You have been defeated!")
		}
	} else {
		cs.ChangeCombatState(types.PlayerTurn)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.PlayerTurn)
		}
	}
}

func (cs *CombatSystem) PlayerAttack(e *entities.Enemy, p *types.Player) {
	damage := p.CalculateDamage(e.Defense)
	message := fmt.Sprintf("%s attacks %s for %d damage!", p.Name, e.Name, damage)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(p.Name, "Attack", e.Name, damage, message)
	}
	
	defeated := e.TakeDamage(damage)
	if defeated {
		cs.ChangeCombatState(types.Victory)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.Victory)
		}
		p.AddExperience(e.ExpReward)
		expMessage := fmt.Sprintf("%s gains %d experience!", p.Name, e.ExpReward)
		if cs.combatUI != nil {
			cs.combatUI.AddAction("System", "Experience", "", 0, expMessage)
			cs.combatUI.AddAction("System", "Result", "", 0, "Victory!")
		}
	} else {
		cs.ChangeCombatState(types.EnemyTurn)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.EnemyTurn)
		}
	}
}

func (cs *CombatSystem) PlayerDefend(p *types.Player) {
	// Defending reduces incoming damage for the next attack
	message := fmt.Sprintf("%s takes a defensive stance!", p.Name)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(p.Name, "Defend", "", 0, message)
	}
	cs.ChangeCombatState(types.EnemyTurn)
	if cs.combatUI != nil {
		cs.combatUI.SetTurn(types.EnemyTurn)
	}
}

func (cs *CombatSystem) PlayerRun(p *types.Player) bool {
	// Simple run calculation - higher speed increases success chance
	successChance := 50 + (p.Stats.Speed * 2) // Base 50% + 2% per speed point
	if successChance > 90 {
		successChance = 90 // Cap at 90%
	}
	
	if rand.Intn(100) < successChance {
		message := fmt.Sprintf("%s successfully runs away!", p.Name)
		if cs.combatUI != nil {
			cs.combatUI.AddAction(p.Name, "Run", "", 0, message)
		}
		cs.ExitCombat()
		return true
	} else {
		message := fmt.Sprintf("%s failed to run away!", p.Name)
		if cs.combatUI != nil {
			cs.combatUI.AddAction(p.Name, "Run", "", 0, message)
		}
		cs.ChangeCombatState(types.EnemyTurn)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.EnemyTurn)
		}
		return false
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
	if cs.combatUI != nil {
		cs.combatUI.AddAction("System", "Combat", "", 0, "Combat ended.")
	}
}

// GetCombatUI returns the combat UI instance
func (cs *CombatSystem) GetCombatUI() *ui.CombatUI {
	return cs.combatUI
}

// ProcessPlayerAction processes the selected player action
func (cs *CombatSystem) ProcessPlayerAction(action string, p *types.Player) bool {
	if cs.CurrentCombatState != types.PlayerTurn || cs.CurrentEnemy == nil {
		return false
	}
	
	switch action {
	case "Attack":
		cs.PlayerAttack(cs.CurrentEnemy, p)
	case "Defend":
		cs.PlayerDefend(p)
	case "Use Item":
		// TODO: Implement item usage
		message := fmt.Sprintf("%s tries to use an item but has none!", p.Name)
		if cs.combatUI != nil {
			cs.combatUI.AddAction(p.Name, "Use Item", "", 0, message)
		}
		cs.ChangeCombatState(types.EnemyTurn)
		if cs.combatUI != nil {
			cs.combatUI.SetTurn(types.EnemyTurn)
		}
	case "Run":
		return cs.PlayerRun(p)
	default:
		return false
	}
	
	return true
}

// Update should be called each frame to handle AI turns and UI updates
func (cs *CombatSystem) Update(p *types.Player) {
	if cs.CurrentCombatState == types.EnemyTurn && cs.CurrentEnemy != nil {
		cs.AiAttack(*cs.CurrentEnemy, p)
	}
	
	// Update the combat UI display
	if cs.IsInCombat() && cs.combatUI != nil {
		cs.combatUI.Display()
	}
}
