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
		Inventory: make([]Item, 0, 10),
		Implants:  [5]Implant{},
		MaxInv:    10,
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
