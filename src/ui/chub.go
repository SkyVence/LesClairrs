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

type CAction struct {
	Timestamp  time.Time
	Actor      string
	ActionType string
	Target     string
	Damage     int
	Message    string
}

type CHistory struct {
	Actions    []CAction
	MaxActions int
}

type CombatHud struct {
	Width      int
	Height     int
	TermWidth  int
	TermHeight int

	Player      *types.Player
	Enemy       *entities.Enemy
	CurrentTurn types.CombatState
	History     *CHistory
	LocManager  *engine.LocalizationManager

	SelectedAction   int
	AvailableActions []string
	ShowHistory      bool

	Styles CHudStyles
}

type CHudStyles struct {
	Container        lipgloss.Style
	TopHealthBar     lipgloss.Style
	HealthBar        lipgloss.Style
	TopEnemyBar      lipgloss.Style
	EnemyBar         lipgloss.Style
	BorderContainer  lipgloss.Style
	Text             lipgloss.Style
	SelectedAction   lipgloss.Style
	UnselectedAction lipgloss.Style
	History          lipgloss.Style
	Victory          lipgloss.Style
	Defeat           lipgloss.Style
}

func DefaultCHudStyles() CHudStyles {
	return CHudStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2),
		TopHealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffffff")).
			Bold(true),
		HealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3aa136ff")),
		TopEnemyBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e74848ff")).
			Bold(true),
		EnemyBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e03c3cff")),
		BorderContainer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")),
		SelectedAction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true),
		UnselectedAction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")),
		History: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Italic(true),
		Victory: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		Defeat: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
	}
}

func NewCombatHistory(maxActions int) *CHistory {
	return &CHistory{
		Actions:    make([]CAction, 0, maxActions),
		MaxActions: maxActions,
	}
}

// AddAction adds a new action to the combat history
func (ch *CHistory) AddAction(action CAction) {
	ch.Actions = append(ch.Actions, action)

	// Remove oldest actions if we exceed the maximum
	if len(ch.Actions) > ch.MaxActions {
		ch.Actions = ch.Actions[1:]
	}
}

// GetRecentActions returns the most recent actions (up to count)
func (ch *CHistory) GetRecentActions(count int) []CAction {
	if count >= len(ch.Actions) {
		return ch.Actions
	}
	return ch.Actions[len(ch.Actions)-count:]
}

// Clear removes all actions from history
func (ch *CHistory) Clear() {
	ch.Actions = ch.Actions[:0]
}

// NewCombatHud creates a new combat UI instance
func NewCombatHud(renderer engine.Renderer, locManager *engine.LocalizationManager) *CombatHud {
	width, height := renderer.GetSize()

	// Fallback to reasonable defaults if renderer size is 0
	if width == 0 {
		width = 80 // Default terminal width
	}
	if height == 0 {
		height = 24 // Default terminal height
	}

	return &CombatHud{
		Width:            width,
		Height:           height,
		TermWidth:        width,
		TermHeight:       height,
		LocManager:       locManager,
		History:          NewCombatHistory(50), // Keep last 50 actions
		SelectedAction:   0,
		AvailableActions: []string{"Attack", "Defend", "Use Item", "Run"},
		ShowHistory:      false,
		Styles:           DefaultCHudStyles(),
	}
}

// SetCombatants sets the player and enemy for the combat
func (cui *CombatHud) SetCombatants(player *types.Player, enemy *entities.Enemy) {
	cui.Player = player
	cui.Enemy = enemy
	cui.History.Clear()
	cui.SelectedAction = 0
}

// SetTurn updates the current turn state
func (cui *CombatHud) SetTurn(turn types.CombatState) {
	cui.CurrentTurn = turn
}

// AddAction adds a new action to combat history using CAction struct
func (cui *CombatHud) AddAction(actor, actionType, target string, damage int, message string) {
	action := CAction{
		Timestamp:  time.Now(),
		Actor:      actor,
		ActionType: actionType,
		Target:     target,
		Damage:     damage,
		Message:    message,
	}
	cui.History.AddAction(action)
}

// UpdateState updates the combat UI state
func (cui *CombatHud) UpdateState(turn types.CombatState) {
	cui.CurrentTurn = turn

	// Clear enemy reference when player is defeated for immediate UI cleanup
	if turn == types.Dead {
		cui.Enemy = nil
	}
}

func (cui *CombatHud) Init() engine.Cmd {
	return engine.TickNow()
}

func (cui *CombatHud) Update(msg engine.Msg) (CombatHud, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		cui.TermWidth = msg.Width
		cui.TermHeight = msg.Height
		cui.Width = msg.Width
		cui.Height = msg.Height
	}
	return *cui, nil
}

func (cui *CombatHud) PHealthBar(current, max int, barWidth int, player *types.Player) string {
	healthPercent := float64(current) / float64(max)
	healthFilled := int(healthPercent * float64(barWidth))
	if healthFilled < 0 {
		healthFilled = 0
	}
	if healthFilled > barWidth {
		healthFilled = barWidth
	}
	healthBar := strings.Repeat("█", healthFilled) + strings.Repeat("▒", barWidth-healthFilled)
	return healthBar
}

func (cui *CombatHud) EHealthBar(current, max int, barWidth int, enemy *entities.Enemy) string {
	healthPercent := float64(current) / float64(max)
	healthFilled := int(healthPercent * float64(barWidth))
	if healthFilled < 0 {
		healthFilled = 0
	}
	if healthFilled > barWidth {
		healthFilled = barWidth
	}
	healthBar := strings.Repeat("█", healthFilled) + strings.Repeat("▒", barWidth-healthFilled)
	return healthBar
}

func (cui *CombatHud) HistoryView(maxLines int) string {

	// Use full terminal height for history
	availableHeight := cui.TermHeight - 4 // Account for borders and padding
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Text Fields
	title := fmt.Sprintf("%s:", cui.LocManager.Text("ui.hud.history.title"))
	missingActions := fmt.Sprintf("(%s)", cui.LocManager.Text("ui.hud.history.no_actions"))
	recentActions := cui.History.GetRecentActions(availableHeight - 2) // Leave space for title

	content := title
	if len(recentActions) == 0 {
		content += "\n" + missingActions
	} else {
		for _, action := range recentActions {
			timeStr := action.Timestamp.Format("15:04:05")
			line := fmt.Sprintf("[%s] %s", timeStr, action.Message)
			if action.Damage > 0 {
				line += fmt.Sprintf(" (-%d)", action.Damage)
			}
			content += "\n" + line
		}
	}

	// Create a style that takes full height
	historyStyle := cui.Styles.Container.
		Height(availableHeight).
		Width(40) // Fixed width for history panel

	return historyStyle.Render(content)
}

func (cui *CombatHud) ActionMenu() string {
	if cui.CurrentTurn != types.PlayerTurn {
		var turnText string
		switch cui.CurrentTurn {
		case types.EnemyTurn:
			turnText = "Enemy Turn - Waiting..."
		case types.Victory:
			turnText = cui.Styles.Victory.Render("Victory!")
		case types.Dead:
			turnText = cui.Styles.Defeat.Render("Defeat!")
		default:
			turnText = "Processing..."
		}
		return cui.Styles.Text.Render(turnText)
	}

	content := cui.Styles.Text.Render(cui.LocManager.Text("ui.hud.actions.prompt") + "\n")

	for i, action := range cui.AvailableActions {
		if i == cui.SelectedAction {
			content += cui.Styles.SelectedAction.Render("> "+cui.LocManager.Text("ui.hud.actions."+strings.ToLower(strings.ReplaceAll(action, " ", "_")))) + "\n"
		} else {
			content += cui.Styles.UnselectedAction.Render("  "+cui.LocManager.Text("ui.hud.actions."+strings.ToLower(strings.ReplaceAll(action, " ", "_")))) + "\n"
		}
	}

	content += "\n" + cui.Styles.Text.Render(cui.LocManager.Text("ui.hud.actions.navigate"))
	return cui.Styles.Container.Render(content)
}

func (cui *CombatHud) InfoView(playerHealthBar string, enemyHealthBar string) string {

	// Text Fields
	playerHealthText := fmt.Sprintf("%s: %d/%d", cui.LocManager.Text("ui.hud.health"), cui.Player.Stats.CurrentHP, cui.Player.Stats.MaxHP)

	// Create player info box
	playerContent := fmt.Sprintf("%s\n%s\n%s",
		cui.Styles.TopHealthBar.Render(cui.Player.Name),
		cui.Styles.Text.Render(playerHealthText),
		cui.Styles.HealthBar.Render(playerHealthBar))
	playerBox := cui.Styles.Container.Render(playerContent)

	// Only show enemy info if combat is still active (not victory/defeat), enemy health bar is provided, and enemy exists
	if cui.CurrentTurn == types.Victory || cui.CurrentTurn == types.Dead || enemyHealthBar == "" || cui.Enemy == nil {
		// Just show player info during victory/defeat screen or when no enemy data is provided
		return playerBox
	}

	// Create enemy info box for active combat
	enemyHealthText := fmt.Sprintf("%s: %d/%d", cui.LocManager.Text("ui.hud.health"), cui.Enemy.CurrentHP, cui.Enemy.MaxHP)
	enemyContent := fmt.Sprintf("%s\n%s\n%s",
		cui.Styles.TopEnemyBar.Render(cui.Enemy.Name),
		cui.Styles.Text.Render(enemyHealthText),
		cui.Styles.EnemyBar.Render(enemyHealthBar))
	enemyBox := cui.Styles.Container.Render(enemyContent)

	// Combine both boxes horizontally for compact layout
	return lipgloss.JoinHorizontal(lipgloss.Top, playerBox, enemyBox)
}

func (cui *CombatHud) View() string {
	if cui.TermWidth == 0 || cui.TermHeight == 0 {
		return "If you see this, something went wrong. The terminal size is zero."
	}

	// Always calculate player health bar
	playerHealthBar := cui.PHealthBar(cui.Player.Stats.CurrentHP, cui.Player.Stats.MaxHP, 20, cui.Player)

	// Only calculate enemy health bar if combat is still active (not victory/defeat) and enemy exists
	var enemyHealthBar string
	if cui.CurrentTurn != types.Victory && cui.CurrentTurn != types.Dead && cui.Enemy != nil {
		enemyHealthBar = cui.EHealthBar(cui.Enemy.CurrentHP, cui.Enemy.MaxHP, 20, cui.Enemy)
	}

	// Get UI components
	infoView := cui.InfoView(playerHealthBar, enemyHealthBar)
	actionMenu := cui.ActionMenu()
	history := cui.HistoryView(cui.History.MaxActions)

	// Create left side with info and action menu stacked vertically
	leftSection := lipgloss.JoinVertical(lipgloss.Left, infoView, actionMenu)

	// Create main layout with left section and history on the right
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftSection, history)
	centeredLayout := lipgloss.NewStyle().Width(cui.TermWidth).Align(lipgloss.Center).Render(mainLayout)

	// Add victory overlay if player won
	if cui.CurrentTurn == types.Victory {
		victoryText := cui.Styles.Victory.Render("🎉 VICTORY! 🎉")
		victoryBanner := lipgloss.NewStyle().
			Width(cui.TermWidth).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#000000")).
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0).
			Render(victoryText)

		// Combine main layout with victory banner at the top
		return lipgloss.JoinVertical(lipgloss.Left, victoryBanner, centeredLayout)
	}

	// Add defeat overlay if player lost
	if cui.CurrentTurn == types.Dead {
		defeatText := cui.Styles.Defeat.Render("💀 DEFEAT 💀")
		defeatBanner := lipgloss.NewStyle().
			Width(cui.TermWidth).
			Align(lipgloss.Center).
			Background(lipgloss.Color("#000000")).
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Padding(1, 0).
			Render(defeatText)

		// Combine main layout with defeat banner at the top
		return lipgloss.JoinVertical(lipgloss.Left, defeatBanner, centeredLayout)
	}

	// Return the centered layout
	return centeredLayout
}
