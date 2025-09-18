package systems

import (
	"projectred-rpg.com/game/types"
)

// NPCSystem manages NPCs and their interactions
type NPCSystem struct {
	npcs        map[string]*types.NPC
	dialogSys   *DialogSystem
	interaction *types.NPC // Currently interacting NPC
}

// NewNPCSystem creates a new NPC management system
func NewNPCSystem(dialogSystem *DialogSystem) *NPCSystem {
	return &NPCSystem{
		npcs:      make(map[string]*types.NPC),
		dialogSys: dialogSystem,
	}
}

// AddNPC adds an NPC to the system
func (ns *NPCSystem) AddNPC(npc *types.NPC) {
	ns.npcs[npc.ID] = npc
}

// RemoveNPC removes an NPC from the system
func (ns *NPCSystem) RemoveNPC(id string) {
	delete(ns.npcs, id)
}

// GetNPC returns an NPC by ID
func (ns *NPCSystem) GetNPC(id string) *types.NPC {
	return ns.npcs[id]
}

// GetAllNPCs returns all NPCs
func (ns *NPCSystem) GetAllNPCs() map[string]*types.NPC {
	return ns.npcs
}

// CheckInteractions checks if the player can interact with any nearby NPCs
func (ns *NPCSystem) CheckInteractions(playerPos types.Position, maxDistance int) *types.NPC {
	for _, npc := range ns.npcs {
		if npc.IsInRange(playerPos, maxDistance) {
			return npc
		}
	}
	return nil
}

// StartInteraction starts a dialog with the specified NPC
func (ns *NPCSystem) StartInteraction(npc *types.NPC, playerName string) {
	if npc == nil || !npc.IsActive {
		return
	}

	ns.interaction = npc
	
	// Create a simple greeting dialog
	dialogKey := npc.GetDialogKey() + ".greeting"
	dialog := CreateSimpleDialog("", dialogKey, playerName)
	
	ns.dialogSys.StartDialog(dialog, npc.Pos, func() {
		ns.EndInteraction()
	})
}

// StartCustomDialog starts a custom dialog sequence with an NPC
func (ns *NPCSystem) StartCustomDialog(npc *types.NPC, dialog *DialogSequence) {
	if npc == nil || !npc.IsActive {
		return
	}

	ns.interaction = npc
	ns.dialogSys.StartDialog(dialog, npc.Pos, func() {
		ns.EndInteraction()
	})
}

// EndInteraction ends the current NPC interaction
func (ns *NPCSystem) EndInteraction() {
	ns.interaction = nil
}

// GetCurrentInteraction returns the currently interacting NPC
func (ns *NPCSystem) GetCurrentInteraction() *types.NPC {
	return ns.interaction
}

// IsInteracting returns whether the player is currently interacting with an NPC
func (ns *NPCSystem) IsInteracting() bool {
	return ns.interaction != nil && ns.dialogSys.IsActive()
}

// CreateGreetingDialog creates a simple greeting dialog for an NPC
func (ns *NPCSystem) CreateGreetingDialog(npc *types.NPC, playerName string) *DialogSequence {
	dialogKey := npc.GetDialogKey() + ".greeting"
	speakerKey := npc.GetDialogKey() + ".name"
	return CreateSimpleDialog(speakerKey, dialogKey, playerName)
}

// CreateMultiPartDialog creates a multi-part dialog for an NPC
func (ns *NPCSystem) CreateMultiPartDialog(npc *types.NPC, dialogKeys []string, playerName string) *DialogSequence {
	if len(dialogKeys) == 0 {
		return nil
	}

	speakerKey := npc.GetDialogKey() + ".name"
	baseKey := npc.GetDialogKey() + "."
	
	entries := make([]DialogEntry, len(dialogKeys))
	for i, key := range dialogKeys {
		entries[i] = DialogEntry{
			SpeakerKey: speakerKey,
			TextKey:    baseKey + key,
			Args:       []any{playerName},
		}
	}
	
	return CreateMultiDialog(entries...)
}

// Example usage functions for different dialog types:

// StartMerchantDialog starts a merchant interaction dialog
func (ns *NPCSystem) StartMerchantDialog(npc *types.NPC, playerName string) {
	if npc.Type != types.NPCMerchant {
		return
	}
	
	dialog := ns.CreateMultiPartDialog(npc, []string{"greeting", "purchase", "farewell"}, playerName)
	ns.StartCustomDialog(npc, dialog)
}

// StartQuestDialog starts a quest-giving dialog
func (ns *NPCSystem) StartQuestDialog(npc *types.NPC, playerName string) {
	if npc.Type != types.NPCAethelgard {
		return
	}
	
	dialog := ns.CreateMultiPartDialog(npc, []string{"greeting", "quest", "farewell"}, playerName)
	ns.StartCustomDialog(npc, dialog)
}

// CreateNPCSprite creates a simple sprite for an NPC (placeholder)
func CreateNPCSprite(npcType types.NPCType) string {
	switch npcType {
	case types.NPCGuard:
		return "🛡️\n👤\n🦵"
	case types.NPCMerchant:
		return "💰\n👤\n🦵"
	case types.NPCVillager:
		return "👒\n👤\n🦵"
	case types.NPCAethelgard:
		return "🎩\n👤\n🦵"
	default:
		return "?\n👤\n🦵"
	}
}