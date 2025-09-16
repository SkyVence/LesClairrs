package aigen

// Backed up player.go from the experimental game renderer
// This file contains the exact content that was in src/player.go before it is removed from src.

// Player represents the stick man character
type Player struct {
    X, Y   int    // Position on screen
    Width  int    // Screen width for boundary checking
    Height int    // Screen height for boundary checking
    Sprite string // ASCII representation of the stick man
}

// NewPlayer creates a new stick man player
func NewPlayer() *Player {
    return &Player{
        X:      5, // Starting position (away from border)
        Y:      3,
        Sprite: createStickManSprite(),
    }
}

// createStickManSprite returns a simple ASCII stick man
func createStickManSprite() string {
    return ` o
/|\
/ \`
}

// Move updates the player position based on direction
func (p *Player) Move(direction rune) {
    switch direction {
    case '↑':
        if p.Y > 1 { // Account for top border
            p.Y--
        }
    case '↓':
        if p.Y < p.Height-4 { // Account for sprite height and bottom border
            p.Y++
        }
    case '←':
        if p.X > 1 { // Account for left border
            p.X--
        }
    case '→':
        if p.X < p.Width-4 { // Account for sprite width and right border
            p.X++
        }
    }
}

// SetBounds sets the screen boundaries for the player
func (p *Player) SetBounds(width, height int) {
    p.Width = width
    p.Height = height
}

// Render returns the stick man positioned on the screen
func (p *Player) Render() string {
    // Create a simple positioned sprite representation
    // This is a basic implementation - we'll improve it in the main render loop
    return p.Sprite
}

// GetPosition returns the current player position
func (p *Player) GetPosition() (int, int) {
    return p.X, p.Y
}
package aigen
package aigen

// Backed up player.go from the experimental game renderer
// Original file moved here for safekeeping.

// ...original content archived...
