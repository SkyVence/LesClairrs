package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/entities"
	"projectred-rpg.com/game/types"
)

// CombatAction represents a single action in combat history
type CombatAction struct {
	Timestamp  time.Time
	Actor      string // "Player" or enemy name
	ActionType string // "Attack", "Defend", "Use Item", "Special"
	Target     string // Target of the action
	Damage     int    // Damage dealt, 0 if not applicable
	Message    string // Formatted message for display
}

// CombatHistory manages the history of combat actions
type CombatHistory struct {
	Actions    []CombatAction
	MaxActions int // Maximum number of actions to keep
}

// CombatUI represents the full-screen turn-based combat interface
type CombatUI struct {
	width          int
	height         int
	termWidth      int
	termHeight     int
	renderer       engine.Renderer
	locManager     *engine.LocalizationManager
	
	// Combat state
	player         *types.Player
	enemy          *entities.Enemy
	currentTurn    types.CombatState
	history        *CombatHistory
	
	// UI state
	selectedAction int
	availableActions []string
	showHistory    bool
	
	// Styles
	styles         CombatUIStyles
}

// CombatUIStyles contains all styling for the combat interface
type CombatUIStyles struct {
	Container      lipgloss.Style
	PlayerPanel    lipgloss.Style
	EnemyPanel     lipgloss.Style
	HistoryPanel   lipgloss.Style
	ActionPanel    lipgloss.Style
	HealthBar      lipgloss.Style
	SelectedAction lipgloss.Style
	NormalAction   lipgloss.Style
	Title          lipgloss.Style
	Text           lipgloss.Style
	Border         lipgloss.Style
	Warning        lipgloss.Style
	Success        lipgloss.Style
	Damage         lipgloss.Style
}

// NewCombatHistory creates a new combat history tracker
func NewCombatHistory(maxActions int) *CombatHistory {
	return &CombatHistory{
		Actions:    make([]CombatAction, 0, maxActions),
		MaxActions: maxActions,
	}
}

// AddAction adds a new action to the combat history
func (ch *CombatHistory) AddAction(action CombatAction) {
	ch.Actions = append(ch.Actions, action)
	
	// Remove oldest actions if we exceed the maximum
	if len(ch.Actions) > ch.MaxActions {
		ch.Actions = ch.Actions[1:]
	}
}

// GetRecentActions returns the most recent actions (up to count)
func (ch *CombatHistory) GetRecentActions(count int) []CombatAction {
	if count >= len(ch.Actions) {
		return ch.Actions
	}
	return ch.Actions[len(ch.Actions)-count:]
}

// Clear removes all actions from history
func (ch *CombatHistory) Clear() {
	ch.Actions = ch.Actions[:0]
}

// NewCombatUI creates a new combat UI instance
func NewCombatUI(renderer engine.Renderer, locManager *engine.LocalizationManager) *CombatUI {
	width, height := renderer.GetSize()
	
	return &CombatUI{
		width:          width,
		height:         height,
		termWidth:      width,
		termHeight:     height,
		renderer:       renderer,
		locManager:     locManager,
		history:        NewCombatHistory(50), // Keep last 50 actions
		selectedAction: 0,
		availableActions: []string{"Attack", "Defend", "Use Item", "Run"},
		showHistory:    false,
		styles:         DefaultCombatUIStyles(),
	}
}

// DefaultCombatUIStyles returns the default combat UI styling
func DefaultCombatUIStyles() CombatUIStyles {
	return CombatUIStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1),
		PlayerPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1).
			Width(30),
		EnemyPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1).
			Width(30),
		HistoryPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFA726")).
			Padding(1),
		ActionPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#66BB6A")).
			Padding(1),
		HealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
		SelectedAction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true).
			Padding(0, 1),
		NormalAction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true).
			Align(lipgloss.Center),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true),
		Damage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true),
	}
}

// UpdateSize updates the UI dimensions when terminal is resized
func (cui *CombatUI) UpdateSize() {
	cui.termWidth, cui.termHeight = cui.renderer.GetSize()
	cui.width = cui.termWidth
	cui.height = cui.termHeight
}

// SetCombatants sets the player and enemy for the combat
func (cui *CombatUI) SetCombatants(player *types.Player, enemy *entities.Enemy) {
	cui.player = player
	cui.enemy = enemy
	cui.history.Clear()
	cui.selectedAction = 0
}

// SetTurn updates the current turn state
func (cui *CombatUI) SetTurn(turn types.CombatState) {
	cui.currentTurn = turn
}

// ToggleHistory toggles the visibility of the combat history
func (cui *CombatUI) ToggleHistory() {
	cui.showHistory = !cui.showHistory
}

// SelectNextAction moves to the next available action
func (cui *CombatUI) SelectNextAction() {
	cui.selectedAction = (cui.selectedAction + 1) % len(cui.availableActions)
}

// SelectPrevAction moves to the previous available action
func (cui *CombatUI) SelectPrevAction() {
	cui.selectedAction--
	if cui.selectedAction < 0 {
		cui.selectedAction = len(cui.availableActions) - 1
	}
}

// GetSelectedAction returns the currently selected action
func (cui *CombatUI) GetSelectedAction() string {
	if cui.selectedAction >= 0 && cui.selectedAction < len(cui.availableActions) {
		return cui.availableActions[cui.selectedAction]
	}
	return ""
}

// AddAction adds an action to the combat history
func (cui *CombatUI) AddAction(actor, actionType, target string, damage int, message string) {
	action := CombatAction{
		Timestamp:  time.Now(),
		Actor:      actor,
		ActionType: actionType,
		Target:     target,
		Damage:     damage,
		Message:    message,
	}
	cui.history.AddAction(action)
}

// renderHealthBar creates a visual health bar
func (cui *CombatUI) renderHealthBar(current, max int, width int) string {
	if max <= 0 {
		return ""
	}
	
	percentage := float64(current) / float64(max)
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 1 {
		percentage = 1
	}
	
	filledWidth := int(float64(width) * percentage)
	emptyWidth := width - filledWidth
	
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)
	
	var color lipgloss.Color
	if percentage > 0.6 {
		color = lipgloss.Color("#4ECDC4") // Green
	} else if percentage > 0.3 {
		color = lipgloss.Color("#FFA726") // Orange
	} else {
		color = lipgloss.Color("#FF6B6B") // Red
	}
	
	bar := lipgloss.NewStyle().Foreground(color).Render(filled) + 
		  lipgloss.NewStyle().Foreground(lipgloss.Color("#333333")).Render(empty)
	
	return fmt.Sprintf("%s %d/%d", bar, current, max)
}

// renderPlayerPanel creates the player information panel
func (cui *CombatUI) renderPlayerPanel() string {
	if cui.player == nil {
		return cui.styles.PlayerPanel.Render("No Player Data")
	}
	
	content := fmt.Sprintf("%s\n", cui.styles.Title.Render(cui.player.Name))
	content += fmt.Sprintf("Level: %d\n", cui.player.Stats.Level)
	content += fmt.Sprintf("HP: %s\n", 
		cui.renderHealthBar(cui.player.Stats.CurrentHP, cui.player.Stats.MaxHP, 20))
	content += fmt.Sprintf("ATK: %d\n", cui.player.Stats.Force)
	content += fmt.Sprintf("DEF: %d\n", cui.player.Stats.Defense)
	content += fmt.Sprintf("SPD: %d\n", cui.player.Stats.Speed)
	content += fmt.Sprintf("ACC: %d", cui.player.Stats.Accuracy)
	
	return cui.styles.PlayerPanel.Render(content)
}

// renderEnemyPanel creates the enemy information panel
func (cui *CombatUI) renderEnemyPanel() string {
	if cui.enemy == nil {
		return cui.styles.EnemyPanel.Render("No Enemy Data")
	}
	
	content := fmt.Sprintf("%s\n", cui.styles.Title.Render(cui.enemy.Name))
	content += fmt.Sprintf("HP: %s\n", 
		cui.renderHealthBar(cui.enemy.CurrentHP, cui.enemy.MaxHP, 20))
	content += fmt.Sprintf("ATK: %d\n", cui.enemy.Force)
	content += fmt.Sprintf("DEF: %d\n", cui.enemy.Defense)
	content += fmt.Sprintf("SPD: %d\n", cui.enemy.Speed)
	content += fmt.Sprintf("ACC: %d", cui.enemy.Accuracy)
	
	return cui.styles.EnemyPanel.Render(content)
}

// renderActionPanel creates the action selection panel
func (cui *CombatUI) renderActionPanel() string {
	if cui.currentTurn != types.PlayerTurn {
		var turnText string
		switch cui.currentTurn {
		case types.EnemyTurn:
			turnText = "Enemy Turn - Waiting..."
		case types.Victory:
			turnText = cui.styles.Success.Render("Victory!")
		case types.Dead:
			turnText = cui.styles.Warning.Render("Defeat!")
		default:
			turnText = "Processing..."
		}
		return cui.styles.ActionPanel.Render(turnText)
	}
	
	content := cui.styles.Title.Render("Choose Action") + "\n\n"
	
	for i, action := range cui.availableActions {
		if i == cui.selectedAction {
			content += cui.styles.SelectedAction.Render(fmt.Sprintf("> %s", action)) + "\n"
		} else {
			content += cui.styles.NormalAction.Render(fmt.Sprintf("  %s", action)) + "\n"
		}
	}
	
	content += "\n" + cui.styles.Text.Render("Use ↑↓ to select, Enter to confirm")
	content += "\n" + cui.styles.Text.Render("Press H to toggle history")
	
	return cui.styles.ActionPanel.Render(content)
}

// renderHistoryPanel creates the combat history panel
func (cui *CombatUI) renderHistoryPanel(maxLines int) string {
	if !cui.showHistory {
		return ""
	}
	
	content := cui.styles.Title.Render("Combat History") + "\n\n"
	
	recentActions := cui.history.GetRecentActions(maxLines - 2)
	if len(recentActions) == 0 {
		content += cui.styles.Text.Render("No actions yet...")
	} else {
		for _, action := range recentActions {
			timeStr := action.Timestamp.Format("15:04:05")
			line := fmt.Sprintf("[%s] %s", timeStr, action.Message)
			if action.Damage > 0 {
				line = cui.styles.Damage.Render(line)
			} else {
				line = cui.styles.Text.Render(line)
			}
			content += line + "\n"
		}
	}
	
	return cui.styles.HistoryPanel.Render(content)
}

// Render creates the complete combat UI
func (cui *CombatUI) Render() string {
	cui.UpdateSize()
	
	// Calculate available space
	minWidth := 80  // Minimum width for proper display
	minHeight := 20 // Minimum height for proper display
	
	if cui.termWidth < minWidth || cui.termHeight < minHeight {
		// Small screen mode - simplified layout
		return cui.renderSmallScreen()
	}
	
	// Full screen mode
	return cui.renderFullScreen()
}

// renderSmallScreen renders a simplified layout for small terminals
func (cui *CombatUI) renderSmallScreen() string {
	content := cui.styles.Title.Render("COMBAT") + "\n\n"
	
	// Player stats (compact)
	if cui.player != nil {
		content += fmt.Sprintf("You: %s (HP: %d/%d)\n", 
			cui.player.Name, cui.player.Stats.CurrentHP, cui.player.Stats.MaxHP)
	}
	
	// Enemy stats (compact)
	if cui.enemy != nil {
		content += fmt.Sprintf("Enemy: %s (HP: %d/%d)\n\n", 
			cui.enemy.Name, cui.enemy.CurrentHP, cui.enemy.MaxHP)
	}
	
	// Recent history (last 3 actions)
	recentActions := cui.history.GetRecentActions(3)
	for _, action := range recentActions {
		content += fmt.Sprintf("• %s\n", action.Message)
	}
	
	content += "\n"
	
	// Actions
	if cui.currentTurn == types.PlayerTurn {
		content += "Actions:\n"
		for i, action := range cui.availableActions {
			if i == cui.selectedAction {
				content += fmt.Sprintf("> %s\n", action)
			} else {
				content += fmt.Sprintf("  %s\n", action)
			}
		}
	} else {
		content += "Enemy turn...\n"
	}
	
	return cui.styles.Container.Render(content)
}

// renderFullScreen renders the full layout for larger terminals
func (cui *CombatUI) renderFullScreen() string {
	// Top section: Player and Enemy panels side by side
	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cui.renderPlayerPanel(),
		strings.Repeat(" ", 4), // Spacer
		cui.renderEnemyPanel(),
	)
	
	// Calculate remaining space for middle section
	historyHeight := cui.termHeight - 15 // Reserve space for top and bottom sections
	if historyHeight < 5 {
		historyHeight = 5
	}
	
	// Middle section: History (if shown) or spacer
	var middleSection string
	if cui.showHistory {
		middleSection = cui.renderHistoryPanel(historyHeight)
	} else {
		// Show a simplified recent action summary
		recentActions := cui.history.GetRecentActions(3)
		content := cui.styles.Title.Render("Recent Actions") + "\n\n"
		for _, action := range recentActions {
			content += cui.styles.Text.Render("• " + action.Message) + "\n"
		}
		if len(recentActions) == 0 {
			content += cui.styles.Text.Render("No actions yet...")
		}
		middleSection = cui.styles.HistoryPanel.Render(content)
	}
	
	// Bottom section: Action panel
	bottomSection := cui.renderActionPanel()
	
	// Combine all sections
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		"\n",
		middleSection,
		"\n",
		bottomSection,
	)
}

// Clear clears the screen and positions cursor at top
func (cui *CombatUI) Clear() {
	cui.renderer.ClearScreen()
	cui.renderer.SetCursor(0, 0)
}

// Display renders and displays the combat UI
func (cui *CombatUI) Display() {
	cui.Clear()
	content := cui.Render()
	cui.renderer.Write(content)
}
