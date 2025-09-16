package game

func NewPlayer(name string, class Class, pos Position) *Player {
	stats := PlayerStats{
		Level:        1,
		Exp:          0,
		NextLevelExp: 100,
		Force:        class.Force,
		Speed:        class.Speed,
		Defense:      class.Defense,
		Accuracy:     class.Accuracy,
		MaxHP:        class.MaxHP,
		CurrentHP:    class.MaxHP,
	}
	return &Player{
		Name:      name,
		Class:     class,
		Stats:     stats,
		Pos:       pos,
		sprite:    createStickManSprite(),
		Inventory: make([]Item, 0, 10),
		Implants:  [5]Implant{},
		MaxInv:    10,
	}
}

// Placeholder until i implement moving animation
func createStickManSprite() string {
	return ` o
/|\
/ \`
}

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

func (p *Player) AddItemToInventory(item Item) bool {
	if len(p.Inventory) >= p.MaxInv {
		return false // Inventory full
	}
	p.Inventory = append(p.Inventory, item)
	return true
}

func (p *Player) RemoveItemFromInventory(index int) bool {
	if index < 0 || index >= len(p.Inventory) {
		return false // Invalid index
	}
	p.Inventory = append(p.Inventory[:index], p.Inventory[index+1:]...)
	return true
}

func (p *Player) GetPosition() (int, int) {
	return p.Pos.X, p.Pos.Y
}
