// Package ui provides user interface components for the game
//
// CombatHUD Usage Example:
//
//	// Initialize the combat HUD
//	combatHUD := ui.NewCombatHUD(80, 24, player, locManager)
//
//	// Update combat state when entering combat
//	combatHUD.UpdateCombatState(types.PlayerTurn, enemy)
//
//	// In your game loop, handle input and render
//	action, handled := combatHUD.HandleInput(keyPress)
//	if handled && action != nil {
//		switch action.Type {
//		case 0: // Attack
//			combatSystem.PlayerAttack(enemy, player)
//		case 1: // Defend
//			combatSystem.ChangeCombatState(types.Defending)
//		case 2: // Item
//			// Handle item usage
//		case 3: // Flee
//			combatSystem.ChangeCombatState(types.OutOfCombat)
//		}
//	}
//
//	// Render the HUD
//	hudDisplay := combatHUD.Render()
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

// CombatState interface defines the methods needed from the combat system
type CombatState interface {
	IsInCombat() bool
	GetCurrentCombatState() types.CombatState
	GetCurrentEnemy() *entities.Enemy
}

// CombatHUD represents the heads-up display during combat
type CombatHUD struct {
	width      int
	height     int
	termWidth  int
	termHeight int

	// Combat state
	combatState    CombatState
	currentState   types.CombatState
	player         *types.Player
	currentEnemy   *entities.Enemy
	combatLog      []string
	maxLogEntries  int
	selectedAction int // 0: Attack, 1: Defend, 2: Item, 3: Flee

	// UI State
	locManager *engine.LocalizationManager
	styles     CombatHUDStyles
}

// CombatHUDStyles contains styling for the combat HUD
type CombatHUDStyles struct {
	Container      lipgloss.Style
	TurnIndicator  lipgloss.Style
	PlayerPanel    lipgloss.Style
	EnemyPanel     lipgloss.Style
	HealthBar      lipgloss.Style
	HealthBarLow   lipgloss.Style
	ActionButton   lipgloss.Style
	SelectedAction lipgloss.Style
	CombatLog      lipgloss.Style
	LogEntry       lipgloss.Style
	Border         lipgloss.Style
	Text           lipgloss.Style
}

// NewCombatHUD creates a new combat HUD instance
func NewCombatHUD(width, height int, player *types.Player, locManager *engine.LocalizationManager) *CombatHUD {
	return &CombatHUD{
		width:          width,
		height:         height,
		currentState:   types.OutOfCombat,
		player:         player,
		combatLog:      make([]string, 0),
		maxLogEntries:  5,
		selectedAction: 0,
		locManager:     locManager,
		styles:         DefaultCombatHUDStyles(),
	}
}

// DefaultCombatHUDStyles returns the default combat HUD styling
func DefaultCombatHUDStyles() CombatHUDStyles {
	return CombatHUDStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2),
		TurnIndicator: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			Align(lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(0, 1),
		PlayerPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(0, 1).
			Width(25),
		EnemyPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF0000")).
			Padding(0, 1).
			Width(25),
		HealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		HealthBarLow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		ActionButton: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#808080")).
			Padding(0, 1).
			Margin(0, 1),
		SelectedAction: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1).
			Margin(0, 1).
			Bold(true),
		CombatLog: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Height(8).
			Width(50).
			Padding(0, 1),
		LogEntry: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
	}
}

// UpdateSize updates the HUD dimensions
func (chud *CombatHUD) UpdateSize(width, height int) {
	chud.termWidth = width
	chud.termHeight = height
}

// Update updates the combat HUD state
func (chud *CombatHUD) Update(enemy *entities.Enemy) {
	chud.currentEnemy = enemy
}

// UpdateCombatState updates the current combat state
func (chud *CombatHUD) UpdateCombatState(state types.CombatState, enemy *entities.Enemy) {
	chud.currentState = state
	chud.currentEnemy = enemy
}

// IsInCombat returns whether currently in combat
func (chud *CombatHUD) IsInCombat() bool {
	return chud.currentState != types.OutOfCombat && chud.currentState != types.Dead
}

// AddLogEntry adds a new entry to the combat log
func (chud *CombatHUD) AddLogEntry(message string) {
	chud.combatLog = append(chud.combatLog, message)
	if len(chud.combatLog) > chud.maxLogEntries {
		chud.combatLog = chud.combatLog[1:]
	}
}

// ClearLog clears the combat log
func (chud *CombatHUD) ClearLog() {
	chud.combatLog = make([]string, 0)
}

// SetSelectedAction sets the currently selected action
func (chud *CombatHUD) SetSelectedAction(action int) {
	if action >= 0 && action <= 3 {
		chud.selectedAction = action
	}
}

// GetSelectedAction returns the currently selected action
func (chud *CombatHUD) GetSelectedAction() int {
	return chud.selectedAction
}

// NavigateActions handles navigation between action buttons
func (chud *CombatHUD) NavigateActions(direction rune) {
	switch direction {
	case '←', 'a':
		if chud.selectedAction > 0 {
			chud.selectedAction--
		}
	case '→', 'd':
		if chud.selectedAction < 3 {
			chud.selectedAction++
		}
	}
}

// renderTurnIndicator renders the turn indicator
func (chud *CombatHUD) renderTurnIndicator() string {
	var turnText string
	switch chud.currentState {
	case types.PlayerTurn:
		turnText = chud.locManager.Text("your_turn")
	case types.EnemyTurn:
		turnText = chud.locManager.Text("enemy_turn")
	case types.Victory:
		turnText = chud.locManager.Text("victory")
	case types.Dead:
		turnText = chud.locManager.Text("defeat")
	default:
		turnText = chud.locManager.Text("combat")
	}

	return chud.styles.TurnIndicator.Render(turnText)
}

// renderPlayerPanel renders the player information panel
func (chud *CombatHUD) renderPlayerPanel() string {
	if chud.player == nil {
		return chud.styles.PlayerPanel.Render("Player: N/A")
	}

	healthPercent := float64(chud.player.Stats.CurrentHP) / float64(chud.player.Stats.MaxHP)
	healthStyle := chud.styles.HealthBar
	if healthPercent < 0.3 {
		healthStyle = chud.styles.HealthBarLow
	}

	healthBar := chud.renderHealthBar(chud.player.Stats.CurrentHP, chud.player.Stats.MaxHP, 20)

	content := fmt.Sprintf("%s\nLv.%d\n%s\nHP: %s",
		chud.player.Name,
		chud.player.Stats.Level,
		healthStyle.Render(healthBar),
		healthStyle.Render(fmt.Sprintf("%d/%d", chud.player.Stats.CurrentHP, chud.player.Stats.MaxHP)))

	return chud.styles.PlayerPanel.Render(content)
}

// renderEnemyPanel renders the enemy information panel
func (chud *CombatHUD) renderEnemyPanel() string {
	if chud.currentEnemy == nil {
		return chud.styles.EnemyPanel.Render("Enemy: N/A")
	}

	healthPercent := float64(chud.currentEnemy.CurrentHP) / float64(chud.currentEnemy.MaxHP)
	healthStyle := chud.styles.HealthBar
	if healthPercent < 0.3 {
		healthStyle = chud.styles.HealthBarLow
	}

	healthBar := chud.renderHealthBar(chud.currentEnemy.CurrentHP, chud.currentEnemy.MaxHP, 20)

	content := fmt.Sprintf("%s\n%s\nHP: %s",
		chud.currentEnemy.Name,
		healthStyle.Render(healthBar),
		healthStyle.Render(fmt.Sprintf("%d/%d", chud.currentEnemy.CurrentHP, chud.currentEnemy.MaxHP)))

	return chud.styles.EnemyPanel.Render(content)
}

// renderHealthBar creates a visual health bar
func (chud *CombatHUD) renderHealthBar(current, max, width int) string {
	if max <= 0 {
		return strings.Repeat("░", width)
	}

	filled := int(float64(current) / float64(max) * float64(width))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}

	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// renderActionButtons renders the combat action buttons
func (chud *CombatHUD) renderActionButtons() string {
	actions := []string{
		chud.locManager.Text("attack"),
		chud.locManager.Text("defend"),
		chud.locManager.Text("item"),
		chud.locManager.Text("flee"),
	}

	var buttons []string
	for i, action := range actions {
		style := chud.styles.ActionButton
		if i == chud.selectedAction {
			style = chud.styles.SelectedAction
		}
		buttons = append(buttons, style.Render(action))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
}

// renderCombatLog renders the combat log
func (chud *CombatHUD) renderCombatLog() string {
	if len(chud.combatLog) == 0 {
		return chud.styles.CombatLog.Render(chud.locManager.Text("combat_started"))
	}

	var logLines []string
	for _, entry := range chud.combatLog {
		logLines = append(logLines, chud.styles.LogEntry.Render(entry))
	}

	// Pad with empty lines if needed
	for len(logLines) < chud.maxLogEntries {
		logLines = append(logLines, "")
	}

	content := strings.Join(logLines, "\n")
	return chud.styles.CombatLog.Render(content)
}

// Render renders the complete combat HUD
func (chud *CombatHUD) Render() string {
	if !chud.IsInCombat() {
		return ""
	}

	// Top section: Turn indicator
	turnIndicator := chud.renderTurnIndicator()

	// Middle section: Player and Enemy panels side by side
	playerPanel := chud.renderPlayerPanel()
	enemyPanel := chud.renderEnemyPanel()
	combatInfo := lipgloss.JoinHorizontal(lipgloss.Top, playerPanel, "  ", enemyPanel)

	// Action buttons (only show during player turn)
	var actionSection string
	if chud.currentState == types.PlayerTurn {
		actionButtons := chud.renderActionButtons()
		actionSection = lipgloss.JoinVertical(lipgloss.Left,
			chud.styles.Text.Render("Choose your action:"),
			actionButtons)
	}

	// Combat log
	combatLog := chud.renderCombatLog()

	// Combine all sections
	var sections []string
	sections = append(sections, turnIndicator)
	sections = append(sections, combatInfo)
	if actionSection != "" {
		sections = append(sections, actionSection)
	}
	sections = append(sections, combatLog)

	content := lipgloss.JoinVertical(lipgloss.Center, sections...)
	return chud.styles.Container.Render(content)
}

// HandleInput processes input during combat and returns the selected action
func (chud *CombatHUD) HandleInput(key rune) (*CombatAction, bool) {
	if !chud.IsInCombat() {
		return nil, false
	}

	switch chud.currentState {
	case types.PlayerTurn:
		switch key {
		case '←', 'a', '→', 'd':
			chud.NavigateActions(key)
			return nil, true
		case '\r', ' ': // Enter or Space
			return chud.executeSelectedAction()
		case '1':
			chud.SetSelectedAction(0)
			return chud.executeSelectedAction()
		case '2':
			chud.SetSelectedAction(1)
			return chud.executeSelectedAction()
		case '3':
			chud.SetSelectedAction(2)
			return chud.executeSelectedAction()
		case '4':
			chud.SetSelectedAction(3)
			return chud.executeSelectedAction()
		}
	}

	return nil, false
}

// CombatAction represents an action that can be taken in combat
type CombatAction struct {
	Type   int    // 0: Attack, 1: Defend, 2: Item, 3: Flee
	Target string // Target of the action if applicable
}

// executeSelectedAction returns the action to be executed
func (chud *CombatHUD) executeSelectedAction() (*CombatAction, bool) {
	if chud.currentEnemy == nil || chud.player == nil {
		return nil, false
	}

	action := &CombatAction{Type: chud.selectedAction}

	switch chud.selectedAction {
	case 0: // Attack
		chud.AddLogEntry(fmt.Sprintf("%s attacks %s!", chud.player.Name, chud.currentEnemy.Name))
		return action, true
	case 1: // Defend
		chud.AddLogEntry(fmt.Sprintf("%s takes a defensive stance.", chud.player.Name))
		return action, true
	case 2: // Item
		chud.AddLogEntry(fmt.Sprintf("%s searches for an item...", chud.player.Name))
		// TODO: Implement item selection and usage
		return action, true
	case 3: // Flee
		chud.AddLogEntry(fmt.Sprintf("%s attempts to flee!", chud.player.Name))
		return action, true
	}

	return nil, false
}
