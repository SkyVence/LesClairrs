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
	combatUI            *ui.CombatHud
	enemyTurnDelay      int    // Frames to wait before processing enemy turn
	maxEnemyTurnDelay   int    // Maximum delay for enemy turns
	resultDisplayDelay  int    // Frames to wait before auto-exiting combat after victory/defeat
	maxResultDelay      int    // Maximum delay for result display
	onExitCallback      func() // Callback to refresh game state when exiting combat
}

// NewCombatSystem creates a new combat system instance
func NewCombatSystem(initialState types.CombatState, locManager *engine.LocalizationManager, spawnerSystem *SpawnerSystem) *CombatSystem {
	return &CombatSystem{
		CurrentCombatState:  initialState,
		PreviousCombatState: initialState,
		locManager:          locManager,
		spawnerSystem:       spawnerSystem,
		combatUI:            nil, // Will be initialized later when renderer is available
		enemyTurnDelay:      0,
		maxEnemyTurnDelay:   30, // Wait ~0.5 second at 60fps before enemy acts
		resultDisplayDelay:  0,
		maxResultDelay:      180, // Wait ~3 seconds at 60fps to show result
	}
}

// SetRenderer sets the renderer and initializes the combat UI
func (cs *CombatSystem) SetRenderer(renderer engine.Renderer) {
	cs.combatUI = ui.NewCombatHud(renderer, cs.locManager)
}

// SetExitCallback sets a callback function to be called when combat exits
func (cs *CombatSystem) SetExitCallback(callback func()) {
	cs.onExitCallback = callback
}

func (cs *CombatSystem) ChangeCombatState(newState types.CombatState) {
	cs.PreviousCombatState = cs.CurrentCombatState
	cs.CurrentCombatState = newState

	// Reset enemy turn delay when transitioning to enemy turn
	if newState == types.EnemyTurn {
		cs.enemyTurnDelay = cs.maxEnemyTurnDelay
	}

	// Reset result display delay when transitioning to victory or defeat
	if newState == types.Victory || newState == types.Dead {
		cs.resultDisplayDelay = cs.maxResultDelay
	}
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
		cs.combatUI.UpdateState(types.PlayerTurn)
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
			cs.combatUI.UpdateState(types.Dead)
			cs.combatUI.AddAction("System", "Result", "", 0, "You have been defeated!")
		}
	} else {
		cs.ChangeCombatState(types.PlayerTurn)
		if cs.combatUI != nil {
			cs.combatUI.UpdateState(types.PlayerTurn)
		}
	}
}

// AiDefend makes the enemy take a defensive stance
func (cs *CombatSystem) AiDefend(e *entities.Enemy) {
	// Enemy gains temporary defense boost
	e.Defense += 2

	message := fmt.Sprintf("%s takes a defensive stance!", e.Name)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(e.Name, "Defend", "", 0, message)
	}

	cs.ChangeCombatState(types.PlayerTurn)
	if cs.combatUI != nil {
		cs.combatUI.UpdateState(types.PlayerTurn)
	}
}

// AiSpecialAttack performs a powerful but less accurate attack
func (cs *CombatSystem) AiSpecialAttack(e *entities.Enemy, p *types.Player) {
	// Check if special attack hits (70% accuracy)
	if rand.Intn(100) >= 70 {
		message := fmt.Sprintf("%s attempts a special attack but misses!", e.Name)
		if cs.combatUI != nil {
			cs.combatUI.AddAction(e.Name, "Special Attack", p.Name, 0, message)
		}
	} else {
		// Special attack deals 1.5x damage
		baseDamage := int(float64(e.Force) * 1.5)
		damage := baseDamage - p.Stats.Defense
		if damage < 1 {
			damage = 1
		}

		// Apply damage
		p.Stats.CurrentHP -= damage
		if p.Stats.CurrentHP < 0 {
			p.Stats.CurrentHP = 0
		}

		message := fmt.Sprintf("%s uses special attack on %s for %d damage!", e.Name, p.Name, damage)
		if cs.combatUI != nil {
			cs.combatUI.AddAction(e.Name, "Special Attack", p.Name, damage, message)
		}
	}

	// Check if player is defeated
	if cs.IsPlayerDefeated(p) {
		cs.ChangeCombatState(types.Dead)
		if cs.combatUI != nil {
			cs.combatUI.UpdateState(types.Dead)
			cs.combatUI.AddAction("System", "Result", "", 0, "You have been defeated!")
		}
	} else {
		cs.ChangeCombatState(types.PlayerTurn)
		if cs.combatUI != nil {
			cs.combatUI.UpdateState(types.PlayerTurn)
		}
	}
}

// AiHeal makes the enemy restore some health
func (cs *CombatSystem) AiHeal(e *entities.Enemy) {
	// Heal 15-25% of max HP
	healAmount := int(float64(e.MaxHP) * (0.15 + rand.Float64()*0.1))
	if healAmount < 1 {
		healAmount = 1
	}

	e.CurrentHP += healAmount
	if e.CurrentHP > e.MaxHP {
		e.CurrentHP = e.MaxHP
	}

	message := fmt.Sprintf("%s heals for %d HP!", e.Name, healAmount)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(e.Name, "Heal", "", healAmount, message)
	}

	cs.ChangeCombatState(types.PlayerTurn)
	if cs.combatUI != nil {
		cs.combatUI.UpdateState(types.PlayerTurn)
	}
}

// AiTaunt makes the enemy taunt the player, reducing accuracy
func (cs *CombatSystem) AiTaunt(e *entities.Enemy, p *types.Player) {
	// Temporarily reduce player accuracy
	p.Stats.Accuracy -= 10
	if p.Stats.Accuracy < 10 {
		p.Stats.Accuracy = 10
	}

	message := fmt.Sprintf("%s taunts %s, reducing accuracy!", e.Name, p.Name)
	if cs.combatUI != nil {
		cs.combatUI.AddAction(e.Name, "Taunt", p.Name, 0, message)
	}

	cs.ChangeCombatState(types.PlayerTurn)
	if cs.combatUI != nil {
		cs.combatUI.UpdateState(types.PlayerTurn)
	}
}

// ProcessEnemyTurn handles the enemy's turn with random action selection
func (cs *CombatSystem) ProcessEnemyTurn(p *types.Player) {
	if cs.CurrentEnemy == nil {
		return
	}

	enemy := cs.CurrentEnemy

	// Determine available actions based on enemy state and type
	var availableActions []string

	// Always available actions
	availableActions = append(availableActions, "attack")

	// Conditional actions
	if enemy.CurrentHP < enemy.MaxHP/2 {
		// If enemy is below 50% HP, more likely to heal or defend
		availableActions = append(availableActions, "heal", "defend")
	}

	if enemy.Force > 15 {
		// Stronger enemies can use special attacks
		availableActions = append(availableActions, "special_attack")
	}

	if rand.Intn(100) < 20 {
		// 20% chance to add taunt option
		availableActions = append(availableActions, "taunt")
	}

	// Always add defend as a possibility
	if rand.Intn(100) < 30 {
		availableActions = append(availableActions, "defend")
	}

	// Select random action
	selectedAction := availableActions[rand.Intn(len(availableActions))]

	// Execute the selected action
	switch selectedAction {
	case "attack":
		cs.AiAttack(*enemy, p)
	case "defend":
		cs.AiDefend(enemy)
	case "special_attack":
		cs.AiSpecialAttack(enemy, p)
	case "heal":
		cs.AiHeal(enemy)
	case "taunt":
		cs.AiTaunt(enemy, p)
	default:
		// Fallback to attack
		cs.AiAttack(*enemy, p)
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
			cs.combatUI.UpdateState(types.Victory)
		}

		// Remove the defeated enemy from the game world immediately
		if cs.spawnerSystem != nil {
			cs.spawnerSystem.RemoveDefeatedEnemies()
		}

		// Immediately refresh the game space to remove defeated enemies visually
		if cs.onExitCallback != nil {
			cs.onExitCallback()
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
			cs.combatUI.UpdateState(types.EnemyTurn)
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
		cs.combatUI.UpdateState(types.EnemyTurn)
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
			cs.combatUI.UpdateState(types.EnemyTurn)
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

	// Ensure defeated enemies are cleaned up when exiting combat
	if cs.spawnerSystem != nil {
		cs.spawnerSystem.RemoveDefeatedEnemies()
	}
	// Immediately refresh the game world enemy list
	if cs.onExitCallback != nil {
		cs.onExitCallback()
	}

	// Trigger callback to signal end of combat (for state transition)
}

// GetCombatUI returns the combat UI instance
func (cs *CombatSystem) GetCombatUI() *ui.CombatHud {
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
			cs.combatUI.UpdateState(types.EnemyTurn)
		}
	case "Run":
		return cs.PlayerRun(p)
	default:
		return false
	}

	return true
}

// IsReadyToExit returns true if combat has ended and should transition back to exploration
func (cs *CombatSystem) IsReadyToExit() bool {
	return cs.CurrentCombatState == types.Idle
}

// Update should be called each frame to handle AI turns and UI updates
func (cs *CombatSystem) Update(p *types.Player) {
	if cs.CurrentCombatState == types.EnemyTurn && cs.CurrentEnemy != nil {
		// Countdown the delay before processing enemy turn
		if cs.enemyTurnDelay > 0 {
			cs.enemyTurnDelay--
		} else {
			// Process the enemy turn
			cs.ProcessEnemyTurn(p)
		}
	}

	// Handle result display delay for victory/defeat states
	if cs.CurrentCombatState == types.Victory || cs.CurrentCombatState == types.Dead {
		if cs.resultDisplayDelay > 0 {
			cs.resultDisplayDelay--
		} else {
			// Auto-exit combat after showing result
			cs.ExitCombat()
		}
	}
}
