package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

// CombatHUD represents the combat-specific HUD overlay
type CombatHUD struct {
	width      int
	height     int
	termWidth  int
	termHeight int

	// Combat state
	combatState    types.CombatState
	player         *types.Player
	enemy          *entities.Enemy
	selectedAction int
	combatLog      []string
	maxLogEntries  int

	// Available actions
	actions []CombatAction
}

// CombatAction represents an available action in combat
type CombatAction struct {
	Name        string
	Description string
	Enabled     bool
}

// CombatHUDStyles contains styling for the combat HUD
type CombatHUDStyles struct {
	Container       lipgloss.Style
	TurnIndicator   lipgloss.Style
	PlayerSection   lipgloss.Style
	EnemySection    lipgloss.Style
	HealthBar       lipgloss.Style
	EnemyHealthBar  lipgloss.Style
	ActionSelected  lipgloss.Style
	ActionNormal    lipgloss.Style
	ActionDisabled  lipgloss.Style
	CombatLog       lipgloss.Style
	Border          lipgloss.Style
}

// DefaultCombatHUDStyles returns the default combat HUD styling
func DefaultCombatHUDStyles() CombatHUDStyles {
	return CombatHUDStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FF4444")).
			Padding(1, 2).
			Background(lipgloss.Color("#1A1A1A")),
		TurnIndicator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Align(lipgloss.Center),
		PlayerSection: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(0, 1).
			Width(30),
		EnemySection: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(0, 1).
			Width(30),
		HealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true),
		EnemyHealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
		ActionSelected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFD700")).
			Bold(true).
			Padding(0, 1),
		ActionNormal: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1),
		ActionDisabled: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Background(lipgloss.Color("#222222")).
			Padding(0, 1),
		CombatLog: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(0, 1).
			Width(40).
			Height(6),
		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")),
	}
}

// NewCombatHUD creates a new combat HUD component
func NewCombatHUD() *CombatHUD {
	return &CombatHUD{
		width:         100,
		height:        20,
		selectedAction: 0,
		maxLogEntries: 5,
		combatLog:     make([]string, 0),
		actions: []CombatAction{
			{Name: "Attack", Description: "Attack the enemy", Enabled: true},
			{Name: "Defend", Description: "Defend against attacks", Enabled: true},
			{Name: "Items", Description: "Use an item", Enabled: true},
			{Name: "Run", Description: "Attempt to flee", Enabled: true},
		},
	}
}

func (ch *CombatHUD) Init() engine.Cmd {
	return engine.TickNow()
}

func (ch *CombatHUD) Update(msg engine.Msg) (CombatHUD, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		ch.termWidth = msg.Width
		ch.termHeight = msg.Height
		ch.width = msg.Width - 4 // Account for margins
	}
	return *ch, nil
}

// SetCombatState updates the current combat state and participants
func (ch *CombatHUD) SetCombatState(state types.CombatState, player *types.Player, enemy *entities.Enemy) {
	ch.combatState = state
	ch.player = player
	ch.enemy = enemy
}

// SetSelectedAction updates the currently selected action
func (ch *CombatHUD) SetSelectedAction(index int) {
	if index >= 0 && index < len(ch.actions) {
		ch.selectedAction = index
	}
}

// GetSelectedAction returns the currently selected action
func (ch *CombatHUD) GetSelectedAction() int {
	return ch.selectedAction
}

// AddCombatLogEntry adds a new entry to the combat log
func (ch *CombatHUD) AddCombatLogEntry(message string) {
	ch.combatLog = append(ch.combatLog, message)
	if len(ch.combatLog) > ch.maxLogEntries {
		ch.combatLog = ch.combatLog[1:] // Remove oldest entry
	}
}

// ClearCombatLog clears all combat log entries
func (ch *CombatHUD) ClearCombatLog() {
	ch.combatLog = make([]string, 0)
}

// createHealthBar creates a visual health bar
func (ch *CombatHUD) createHealthBar(current, max int, width int, style lipgloss.Style) string {
	if max <= 0 {
		return style.Render(strings.Repeat("▒", width))
	}
	
	percent := float64(current) / float64(max)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	
	filled := int(float64(width) * percent)
	healthBar := strings.Repeat("█", filled) + strings.Repeat("▒", width-filled)
	return style.Render(healthBar)
}

// renderTurnIndicator renders whose turn it is
func (ch *CombatHUD) renderTurnIndicator() string {
	locManager := engine.GetLocalizationManager()
	styles := DefaultCombatHUDStyles()
	
	var turnText string
	switch ch.combatState {
	case types.PlayerTurn:
		turnText = locManager.Text("ui.combat.player_turn")
	case types.EnemyTurn:
		turnText = locManager.Text("ui.combat.enemy_turn")
	case types.Victory:
		turnText = locManager.Text("ui.combat.victory")
	case types.Dead:
		turnText = locManager.Text("ui.combat.defeat")
	default:
		turnText = locManager.Text("ui.combat.combat_active")
	}
	
	return styles.TurnIndicator.Render(fmt.Sprintf("=== %s ===", turnText))
}

// renderPlayerSection renders the player's combat information
func (ch *CombatHUD) renderPlayerSection() string {
	if ch.player == nil {
		return ""
	}
	
	locManager := engine.GetLocalizationManager()
	styles := DefaultCombatHUDStyles()
	
	playerName := ch.player.Name
	playerHP := fmt.Sprintf("%s: %d/%d", 
		locManager.Text("ui.combat.hp"), 
		ch.player.Stats.CurrentHP, 
		ch.player.Stats.MaxHP)
	
	healthBar := ch.createHealthBar(ch.player.Stats.CurrentHP, ch.player.Stats.MaxHP, 20, styles.HealthBar)
	
	playerStats := fmt.Sprintf("%s: %d | %s: %d | %s: %d", 
		locManager.Text("ui.combat.attack"), ch.player.Stats.Force,
		locManager.Text("ui.combat.defense"), ch.player.Stats.Defense,
		locManager.Text("ui.combat.speed"), ch.player.Stats.Speed)
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(playerName),
		playerHP,
		healthBar,
		playerStats,
	)
	
	return styles.PlayerSection.Render(content)
}

// renderEnemySection renders the enemy's combat information
func (ch *CombatHUD) renderEnemySection() string {
	if ch.enemy == nil {
		return ""
	}
	
	locManager := engine.GetLocalizationManager()
	styles := DefaultCombatHUDStyles()
	
	enemyName := ch.enemy.Name
	enemyHP := fmt.Sprintf("%s: %d/%d", 
		locManager.Text("ui.combat.hp"), 
		ch.enemy.CurrentHP, 
		ch.enemy.MaxHP)
	
	healthBar := ch.createHealthBar(ch.enemy.CurrentHP, ch.enemy.MaxHP, 20, styles.EnemyHealthBar)
	
	enemyStats := fmt.Sprintf("%s: %d | %s: %d | %s: %d", 
		locManager.Text("ui.combat.attack"), ch.enemy.Force,
		locManager.Text("ui.combat.defense"), ch.enemy.Defense,
		locManager.Text("ui.combat.speed"), ch.enemy.Speed)
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(enemyName),
		enemyHP,
		healthBar,
		enemyStats,
	)
	
	return styles.EnemySection.Render(content)
}

// renderActions renders the available combat actions
func (ch *CombatHUD) renderActions() string {
	locManager := engine.GetLocalizationManager()
	styles := DefaultCombatHUDStyles()
	
	var actionButtons []string
	for i, action := range ch.actions {
		actionText := locManager.Text("ui.combat.action." + strings.ToLower(action.Name))
		if strings.HasPrefix(actionText, "⟦") && strings.HasSuffix(actionText, "⟧") {
			actionText = action.Name // Fallback to English if translation missing
		}
		
		var style lipgloss.Style
		if !action.Enabled {
			style = styles.ActionDisabled
		} else if i == ch.selectedAction {
			style = styles.ActionSelected
		} else {
			style = styles.ActionNormal
		}
		
		actionButtons = append(actionButtons, style.Render(actionText))
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Top, actionButtons...)
}

// renderCombatLog renders the combat log
func (ch *CombatHUD) renderCombatLog() string {
	styles := DefaultCombatHUDStyles()
	
	logContent := ""
	if len(ch.combatLog) == 0 {
		logContent = lipgloss.NewStyle().Italic(true).Render("Combat log...")
	} else {
		logContent = strings.Join(ch.combatLog, "\n")
	}
	
	return styles.CombatLog.Render(logContent)
}

// View renders the complete combat HUD
func (ch *CombatHUD) View() string {
	if ch.player == nil || ch.enemy == nil {
		return "" // Don't render if combat data not set
	}
	
	styles := DefaultCombatHUDStyles()
	
	// Create all sections
	turnIndicator := ch.renderTurnIndicator()
	playerSection := ch.renderPlayerSection()
	enemySection := ch.renderEnemySection()
	actions := ch.renderActions()
	combatLog := ch.renderCombatLog()
	
	// Layout the combat info (player vs enemy)
	combatInfo := lipgloss.JoinHorizontal(lipgloss.Top,
		playerSection,
		"  VS  ",
		enemySection,
	)
	
	// Layout actions and log side by side
	bottomSection := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Actions:"),
			actions,
		),
		"  ",
		combatLog,
	)
	
	// Combine all sections vertically
	content := lipgloss.JoinVertical(lipgloss.Center,
		turnIndicator,
		"",
		combatInfo,
		"",
		bottomSection,
	)
	
	return styles.Container.Width(ch.width-4).Render(content)
}

// RenderWithContent positions the combat HUD over main content
func (ch *CombatHUD) RenderWithContent(mainContent string) string {
	if ch.termWidth == 0 || ch.termHeight == 0 {
		return mainContent
	}
	
	combatHUDView := ch.View()
	if combatHUDView == "" {
		return mainContent // No combat HUD to show
	}
	
	// Center the combat HUD over the main content
	return lipgloss.Place(
		ch.termWidth, ch.termHeight,
		lipgloss.Center, lipgloss.Center,
		combatHUDView,
		lipgloss.WithWhitespaceChars("·"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333333")),
	)
}