package systems

import (
	"projectred-rpg.com/game/types"
)

// InventorySystem handles inventory management operations
type InventorySystem struct {
	// Inventory configuration and state
}

// NewInventorySystem creates a new inventory system instance
func NewInventorySystem() *InventorySystem {
	return &InventorySystem{}
}

// AddItem attempts to add an item to the player's inventory
func (is *InventorySystem) AddItem(player *types.Player, item types.Item) bool {
	return player.AddItemToInventory(item)
}

// RemoveItem attempts to remove an item from the player's inventory by index
func (is *InventorySystem) RemoveItem(player *types.Player, index int) bool {
	return player.RemoveItemFromInventory(index)
}

// GetInventoryCount returns the current number of items in the inventory
func (is *InventorySystem) GetInventoryCount(player *types.Player) int {
	return len(player.Inventory)
}

// GetInventorySpace returns the available inventory space
func (is *InventorySystem) GetInventorySpace(player *types.Player) int {
	return player.MaxInv - len(player.Inventory)
}

// IsInventoryFull checks if the player's inventory is full
func (is *InventorySystem) IsInventoryFull(player *types.Player) bool {
	return len(player.Inventory) >= player.MaxInv
}

// FindItemByName searches for an item in the inventory by name
func (is *InventorySystem) FindItemByName(player *types.Player, name string) (types.Item, int, bool) {
	for i, item := range player.Inventory {
		if item.Name == name {
			return item, i, true
		}
	}
	return types.Item{}, -1, false
}