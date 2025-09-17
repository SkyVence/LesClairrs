// Package engine provides improved game loop management with frame-based updates.
package engine

import (
	"time"
)

// UpdateManager handles frame-based updates and render coordination
type UpdateManager struct {
	targetFPS  int
	frameDelta time.Duration
	lastUpdate time.Time

	// Update channels
	inputChan  chan Msg
	updateChan chan UpdateMsg
	renderChan chan RenderMsg

	// Performance tracking
	frameCount int
	fpsCounter time.Time
	currentFPS float64
}

// UpdateMsg represents a frame update message
type UpdateMsg struct {
	DeltaTime   time.Duration // Time since last update
	TotalTime   time.Duration // Total game time
	FrameNumber int           // Current frame number
}

// RenderMsg represents a render request
type RenderMsg struct {
	ForceRender bool // Force render even if no changes
}

// NewUpdateManager creates a new update manager
func NewUpdateManager(targetFPS int) *UpdateManager {
	frameDelta := time.Duration(1000/targetFPS) * time.Millisecond

	return &UpdateManager{
		targetFPS:  targetFPS,
		frameDelta: frameDelta,
		lastUpdate: time.Now(),
		inputChan:  make(chan Msg, 100),
		updateChan: make(chan UpdateMsg, 10),
		renderChan: make(chan RenderMsg, 10),
		fpsCounter: time.Now(),
	}
}

// Start begins the game loop with improved frame management
func (um *UpdateManager) Start(model Model) {
	ticker := time.NewTicker(um.frameDelta)
	defer ticker.Stop()

	gameStart := time.Now()

	for {
		select {
		case <-ticker.C:
			// Fixed timestep update
			now := time.Now()
			delta := now.Sub(um.lastUpdate)
			um.lastUpdate = now

			// Send update message
			updateMsg := UpdateMsg{
				DeltaTime:   delta,
				TotalTime:   now.Sub(gameStart),
				FrameNumber: um.frameCount,
			}

			// Non-blocking update send
			select {
			case um.updateChan <- updateMsg:
			default: // Skip frame if update is behind
			}

			um.frameCount++
			um.updateFPS()

		case inputMsg := <-um.inputChan:
			// Handle input immediately
			newModel, cmd := model.Update(inputMsg)
			model = newModel
			if cmd != nil {
				// Execute command in goroutine
				go func() {
					if msg := cmd(); msg != nil {
						um.inputChan <- msg
					}
				}()
			}

		case updateMsg := <-um.updateChan:
			// Process frame update
			newModel, cmd := model.Update(updateMsg)
			model = newModel
			if cmd != nil {
				go func() {
					if msg := cmd(); msg != nil {
						um.inputChan <- msg
					}
				}()
			}

			// Request render
			select {
			case um.renderChan <- RenderMsg{ForceRender: false}:
			default: // Skip render if falling behind
			}
		}
	}
}

// updateFPS calculates current FPS
func (um *UpdateManager) updateFPS() {
	if time.Since(um.fpsCounter) >= time.Second {
		um.currentFPS = float64(um.frameCount) / time.Since(um.fpsCounter).Seconds()
		um.frameCount = 0
		um.fpsCounter = time.Now()
	}
}

// GetFPS returns current frames per second
func (um *UpdateManager) GetFPS() float64 {
	return um.currentFPS
}
