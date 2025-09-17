// Package game provides enhanced update management for ProjectRed RPG
package game

import (
	"time"

	"projectred-rpg.com/engine"
)

// GameUpdateManager enhances the basic update system with advanced features
type GameUpdateManager struct {
	// Update timing
	lastUpdate time.Time
	deltaTime  time.Duration
	gameTime   time.Duration

	// Update optimization
	needsRender   bool
	skipFrames    int
	maxSkipFrames int

	// State tracking
	previousState gameState
	stateChanged  bool

	// Performance tracking
	updateCount   int
	renderCount   int
	lastFPSUpdate time.Time
	currentFPS    float64
}

// NewGameUpdateManager creates an enhanced update manager
func NewGameUpdateManager() *GameUpdateManager {
	return &GameUpdateManager{
		lastUpdate:    time.Now(),
		maxSkipFrames: 5, // Skip max 5 frames if falling behind
		lastFPSUpdate: time.Now(),
	}
}

// ShouldUpdate determines if the game needs an update based on timing
func (gum *GameUpdateManager) ShouldUpdate() bool {
	now := time.Now()
	gum.deltaTime = now.Sub(gum.lastUpdate)

	// Update at least every 16ms (60 FPS target)
	if gum.deltaTime >= 16*time.Millisecond {
		gum.lastUpdate = now
		gum.gameTime += gum.deltaTime
		return true
	}
	return false
}

// ShouldRender determines if rendering is needed
func (gum *GameUpdateManager) ShouldRender() bool {
	if gum.needsRender {
		gum.needsRender = false
		gum.renderCount++
		return true
	}

	// Force render if we've skipped too many frames
	if gum.skipFrames >= gum.maxSkipFrames {
		gum.skipFrames = 0
		gum.renderCount++
		return true
	}

	gum.skipFrames++
	return false
}

// MarkNeedsRender marks that rendering is required
func (gum *GameUpdateManager) MarkNeedsRender() {
	gum.needsRender = true
}

// TrackStateChange tracks when game state changes
func (gum *GameUpdateManager) TrackStateChange(currentState gameState) {
	if currentState != gum.previousState {
		gum.stateChanged = true
		gum.needsRender = true
		gum.previousState = currentState
	} else {
		gum.stateChanged = false
	}
}

// GetDeltaTime returns the time since last update
func (gum *GameUpdateManager) GetDeltaTime() time.Duration {
	return gum.deltaTime
}

// GetGameTime returns total game time
func (gum *GameUpdateManager) GetGameTime() time.Duration {
	return gum.gameTime
}

// UpdateFPS calculates current FPS
func (gum *GameUpdateManager) UpdateFPS() {
	if time.Since(gum.lastFPSUpdate) >= time.Second {
		gum.currentFPS = float64(gum.renderCount) / time.Since(gum.lastFPSUpdate).Seconds()
		gum.renderCount = 0
		gum.lastFPSUpdate = time.Now()
	}
}

// GetFPS returns current frames per second
func (gum *GameUpdateManager) GetFPS() float64 {
	return gum.currentFPS
}

// Enhanced update method for the model
func (m *model) EnhancedUpdate(msg engine.Msg, updateManager *GameUpdateManager) (engine.Model, engine.Cmd) {
	// Track state changes
	updateManager.TrackStateChange(m.state)

	// Handle different message types
	switch msg := msg.(type) {
	case engine.SizeMsg:
		// Terminal resize - always needs render
		m.width = msg.Width
		m.height = msg.Height
		m.menu, _ = m.menu.Update(msg)
		if m.gameSpace == nil {
			m.gameSpace = NewGameRenderer(msg.Width-1, msg.Height-m.hud.Height()-1)
		} else {
			m.gameSpace.UpdateSize(msg.Width-1, msg.Height-m.hud.Height()-1)
		}
		*m.hud, _ = m.hud.Update(msg)
		m.spinner, _ = m.spinner.Update(msg)
		updateManager.MarkNeedsRender()

	case engine.KeyMsg:
		// Input handling
		handled := m.handleKeyInput(msg.Rune, updateManager)
		if handled {
			updateManager.MarkNeedsRender()
		}

	case engine.TickMsg:
		// Time-based updates
		m.handleTimeUpdate(msg, updateManager)

	case GameUpdateMsg:
		// Custom game update messages
		m.handleGameUpdate(msg, updateManager)
	}

	// Update FPS tracking
	updateManager.UpdateFPS()

	return m, nil
}

// GameUpdateMsg represents a custom game update
type GameUpdateMsg struct {
	DeltaTime time.Duration
	GameTime  time.Duration
}

// handleKeyInput processes keyboard input
func (m *model) handleKeyInput(key rune, updateManager *GameUpdateManager) bool {
	switch key {
	case 'q':
		if m.state == stateGame {
			m.state = stateMenu
			return true
		}
		return false

	case '\r', '\n', ' ': // Enter key
		if m.state == stateMenu {
			selected := m.menu.GetSelected()
			switch selected.Value {
			case "start":
				m.state = stateGame
				return true
			}
		}
		return false

	case '↑', '↓', '←', '→':
		if m.state == stateGame {
			oldX, oldY := m.game.Player.GetPosition()
			m.game.Player.Move(key, m.gameSpace.width, m.gameSpace.height)
			newX, newY := m.game.Player.GetPosition()

			// Only mark for render if position actually changed
			return oldX != newX || oldY != newY
		}
		return false
	}

	// Handle menu navigation
	if m.state == stateMenu {
		m.menu, _ = m.menu.Update(engine.KeyMsg{Rune: key})
		return true
	}

	return false
}

// handleTimeUpdate processes time-based updates
func (m *model) handleTimeUpdate(msg engine.TickMsg, updateManager *GameUpdateManager) {
	switch m.state {
	case stateTransition:
		// Update spinner and transition logic
		var cmd engine.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if cmd != nil {
			// Transition complete
			_ = m.game.Advance()
			m.state = stateGame
			updateManager.MarkNeedsRender()
		}

	case stateGame:
		// Update game systems that need time-based updates
		deltaTime := updateManager.GetDeltaTime()
		m.updateGameSystems(deltaTime)
	}
}

// handleGameUpdate processes custom game updates
func (m *model) handleGameUpdate(msg GameUpdateMsg, updateManager *GameUpdateManager) {
	// Handle custom game logic that needs delta time
	m.updateGameSystems(msg.DeltaTime)
}

// updateGameSystems updates time-dependent game systems
func (m *model) updateGameSystems(deltaTime time.Duration) {
	// Future: Update animations, AI, physics, etc.
	// For now, this is a placeholder for time-based updates

	// Example: Update player animations
	// m.game.Player.UpdateAnimations(deltaTime)

	// Example: Update enemy AI
	// for _, enemy := range m.game.CurrentStage.Enemies {
	//     enemy.UpdateAI(deltaTime, m.game.Player)
	// }
}
