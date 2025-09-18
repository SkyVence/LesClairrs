package types

// NPCType represents different types of NPCs
type NPCType string

const (
	NPCGuard    NPCType = "guard"
	NPCMerchant NPCType = "merchant" 
	NPCVillager NPCType = "villager"
	NPCAethelgard NPCType = "aethelgard"
)

// NPC represents a non-player character
type NPC struct {
	ID       string
	Name     string
	Type     NPCType
	Pos      Position
	Sprite   string
	IsActive bool // Whether the NPC can be interacted with
}

// NewNPC creates a new NPC with the specified parameters
func NewNPC(id, name string, npcType NPCType, pos Position, sprite string) *NPC {
	return &NPC{
		ID:       id,
		Name:     name,
		Type:     npcType,
		Pos:      pos,
		Sprite:   sprite,
		IsActive: true,
	}
}

// GetDialogKey returns the base localization key for this NPC's dialogs
func (n *NPC) GetDialogKey() string {
	return "dialog.npcs." + string(n.Type)
}

// IsInRange checks if the player is within interaction range of the NPC
func (n *NPC) IsInRange(playerPos Position, maxDistance int) bool {
	if !n.IsActive {
		return false
	}
	
	dx := abs(n.Pos.X - playerPos.X)
	dy := abs(n.Pos.Y - playerPos.Y)
	
	return dx <= maxDistance && dy <= maxDistance
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetSprite returns the NPC's sprite
func (n *NPC) GetSprite() string {
	return n.Sprite
}

// SetSprite sets the NPC's sprite
func (n *NPC) SetSprite(sprite string) {
	n.Sprite = sprite
}

// SetActive sets whether the NPC can be interacted with
func (n *NPC) SetActive(active bool) {
	n.IsActive = active
}