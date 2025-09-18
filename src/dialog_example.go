package main

import (
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/systems"
	"projectred-rpg.com/game/types"
)

// This file demonstrates how to integrate the dialog system into your game

// GameWithDialog shows how to integrate the dialog system into your main game structure
type GameWithDialog struct {
	player     *types.Player
	dialogSys  *systems.DialogSystem
	npcSys     *systems.NPCSystem
	gameWidth  int
	gameHeight int
}

// NewGameWithDialog creates a new game instance with dialog system
func NewGameWithDialog() *GameWithDialog {
	// Create dialog system with appropriate width for dialog boxes
	dialogSystem := systems.NewDialogSystem(80 - 10) // Leave some margin

	// Create NPC system
	npcSystem := systems.NewNPCSystem(dialogSystem)

	game := &GameWithDialog{
		dialogSys:  dialogSystem,
		npcSys:     npcSystem,
		gameWidth:  80,
		gameHeight: 24,
	}

	// Initialize some example NPCs
	game.initializeNPCs()

	return game
}

func (g *GameWithDialog) Init() engine.Msg {
	return nil
}

// initializeNPCs sets up example NPCs in the game world
func (g *GameWithDialog) initializeNPCs() {
	// Create a merchant NPC
	merchant := types.NewNPC(
		"merchant_1",
		"Shop Keeper",
		types.NPCMerchant,
		types.Position{X: 10, Y: 5},
		systems.CreateNPCSprite(types.NPCMerchant),
	)
	g.npcSys.AddNPC(merchant)

	// Create a guard NPC
	guard := types.NewNPC(
		"guard_1",
		"Town Guard",
		types.NPCGuard,
		types.Position{X: 15, Y: 8},
		systems.CreateNPCSprite(types.NPCGuard),
	)
	g.npcSys.AddNPC(guard)

	// Create Aethelgard (quest giver)
	aethelgard := types.NewNPC(
		"aethelgard_1",
		"Aethelgard",
		types.NPCAethelgard,
		types.Position{X: 20, Y: 10},
		systems.CreateNPCSprite(types.NPCAethelgard),
	)
	g.npcSys.AddNPC(aethelgard)
}

// Update handles game updates including dialog system
func (g *GameWithDialog) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
	// Update dialog system first (it may consume input)
	if g.dialogSys.IsActive() {
		g.dialogSys.Update(msg)
		return nil, nil // Don't process other input while dialog is active
	}

	// Handle regular game input
	if keyMsg, ok := msg.(engine.KeyMsg); ok {
		switch keyMsg.Rune {
		case '↑', '↓', '←', '→':
			// Move player
			if g.player != nil {
				g.player.Move(keyMsg.Rune, g.gameWidth, g.gameHeight)
				g.checkNPCInteractions()
			}
		case 'e', 'E', '\r': // Interaction key
			g.handleInteraction()
		}
	}

	// Update NPC system
	// (NPCs don't need regular updates in this simple example)
	return g, nil
}

// checkNPCInteractions checks if player is near any NPCs
func (g *GameWithDialog) checkNPCInteractions() {
	if g.player == nil {
		return
	}

	// Check if player is near any NPCs
	nearbyNPC := g.npcSys.CheckInteractions(g.player.Pos, 2) // Within 2 tiles

	// You could show an interaction prompt here
	// For example, display "Press E to interact" near the NPC
	if nearbyNPC != nil {
		// Show interaction hint (this would be integrated with your UI system)
		// fmt.Printf("Press E to interact with %s\n", nearbyNPC.Name)
	}
}

// handleInteraction handles player interaction attempts
func (g *GameWithDialog) handleInteraction() {
	if g.player == nil || g.npcSys.IsInteracting() {
		return
	}

	// Check for nearby NPCs
	nearbyNPC := g.npcSys.CheckInteractions(g.player.Pos, 2)
	if nearbyNPC == nil {
		return
	}

	// Start appropriate dialog based on NPC type
	switch nearbyNPC.Type {
	case types.NPCMerchant:
		g.npcSys.StartMerchantDialog(nearbyNPC, g.player.Name)
	case types.NPCAethelgard:
		g.npcSys.StartQuestDialog(nearbyNPC, g.player.Name)
	default:
		// Start simple greeting dialog
		g.npcSys.StartInteraction(nearbyNPC, g.player.Name)
	}
}

// View renders the game including NPCs and dialog
func (g *GameWithDialog) View() string {
	// This is a simplified render example
	// In your actual game, you'd integrate this with your existing renderer

	var result string

	// Render game world (your existing game rendering code goes here)
	// ...

	// Render NPCs
	for _, npc := range g.npcSys.GetAllNPCs() {
		if npc.IsActive {
			// Add NPC to render at npc.Pos with npc.GetSprite()
			// This would be integrated with your existing rendering system
		}
	}

	// Render dialog box on top if active
	if g.dialogSys.IsActive() {
		dialogRender := g.dialogSys.Render()
		// Position and overlay the dialog on your game view
		result += dialogRender
	}

	return result
}

// SetPlayer sets the player for the game
func (g *GameWithDialog) SetPlayer(player *types.Player) {
	g.player = player
}

// Example of creating custom dialog sequences:

// CreateWelcomeDialog creates a welcome dialog sequence
func CreateWelcomeDialog() *systems.DialogSequence {
	return systems.CreateMultiDialog(
		systems.DialogEntry{
			SpeakerKey: "dialog.npcs.aethelgard.name",
			TextKey:    "dialog.npcs.aethelgard.greeting",
			Args:       []any{"Player"}, // Player name would be passed here
		},
		systems.DialogEntry{
			SpeakerKey: "dialog.npcs.aethelgard.name",
			TextKey:    "dialog.npcs.aethelgard.quest",
			Args:       []any{},
		},
		systems.DialogEntry{
			SpeakerKey: "dialog.npcs.aethelgard.name",
			TextKey:    "dialog.npcs.aethelgard.farewell",
			Args:       []any{},
		},
	)
}

// Example of using the dialog system in a specific game event:
func (g *GameWithDialog) TriggerStoryDialog() {
	// Find Aethelgard NPC
	aethelgard := g.npcSys.GetNPC("aethelgard_1")
	if aethelgard != nil {
		// Create a custom story dialog
		storyDialog := CreateWelcomeDialog()
		g.npcSys.StartCustomDialog(aethelgard, storyDialog)
	}
}
