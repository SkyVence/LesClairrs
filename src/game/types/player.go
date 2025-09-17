// Package types defines all core data structures and their associated methods
// for ProjectRed RPG.
//
// This package contains the fundamental types that represent game entities:
//   - Player: Character data, stats, inventory, and behavior methods
//   - Class: Character class definitions with base stats
//   - PlayerStats: Level, experience, and combat statistics
//   - Implant: Cybernetic enhancements with stat bonuses
//
// All types include their associated methods, following Go best practices
// of defining methods in the same package as the type.
//
// Example usage:
//
//	player := &types.Player{Name: "Sam", Class: someClass}
//	player.Move('↑', width, height)
//	success := player.AddItemToInventory(item)
package types

type Class struct {
	Name        string
	Description string
	MaxHP       int
	Force       int
	Speed       int
	Defense     int
	Accuracy    int
}

type BonusStats struct {
	Force    int
	Speed    int
	Defense  int
	Accuracy int
}

type Implant struct {
	Type        BodyParts
	Rarity      Rarity
	Bonus       BonusStats
	Name        string
	Description string
	Cooldown    int
}

type PlayerStats struct {
	Level        int
	Exp          float32
	NextLevelExp int
	Force        int
	Speed        int
	Defense      int
	Accuracy     int
	MaxHP        int
	CurrentHP    int
}

type Player struct {
	Name  string
	Class Class
	Stats PlayerStats
	Pos   Position

	// Placeholder until i implement moving animation
	sprite string

	Inventory []Item
	Implants  [5]Implant // "tete", "brasD", etc - fixed size array
	MaxInv    int
}

// Player methods - defined in the same package as the type
func (p *Player) Move(direction rune, width, height int) {
	switch direction {
	case '↑':
		if p.Pos.Y > 1 { // Account for top border
			p.Pos.Y--
		}
	case '↓':
		if p.Pos.Y < height-4 { // Account for sprite height and bottom border
			p.Pos.Y++
		}
	case '←':
		if p.Pos.X > 1 { // Account for left border
			p.Pos.X--
		}
	case '→':
		if p.Pos.X < width-4 { // Account for sprite width and right border
			p.Pos.X++
		}
	}
}

func (p *Player) GetSprite() string {
	return p.sprite
}

func (p *Player) SetSprite(sprite string) {
	p.sprite = sprite
}

// CreateStickManSprite returns a default player sprite
func CreateStickManSprite() string {
	return ` o
/|\
/ \`
}

// AddItemToInventory adds an item to the player's inventory if there's space
func (p *Player) AddItemToInventory(item Item) bool {
	if len(p.Inventory) >= p.MaxInv {
		return false // Inventory full
	}
	p.Inventory = append(p.Inventory, item)
	return true
}

// RemoveItemFromInventory removes an item from the player's inventory by index
func (p *Player) RemoveItemFromInventory(index int) bool {
	if index < 0 || index >= len(p.Inventory) {
		return false // Invalid index
	}
	p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	return true
}

// GetPosition returns the player's X and Y coordinates
func (p *Player) GetPosition() (int, int) {
	return p.Pos.X, p.Pos.Y
}
