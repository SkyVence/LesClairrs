package game

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/systems"
)

func (gr *GameRender) refreshMenusAfterLanguageChange() {
	locManager := engine.GetLocalizationManager()
	sizeMsg := engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight}

	gr.mainMenu = InitMainMenu(locManager)
	gr.mainMenu, _ = gr.mainMenu.Update(sizeMsg)

	classes := config.GetDefaultClasses()
	gr.classSelection = InitializeClassSelection(locManager, classes)
	gr.classSelection, _ = gr.classSelection.Update(sizeMsg)

	supportedLanguages, err := locManager.GetSupportedLanguages()
	if err != nil {
		supportedLanguages = []string{"fr"}
	}
	gr.settingsMenu = InitializeSettingsSelection(locManager, supportedLanguages)
	gr.settingsMenu, _ = gr.settingsMenu.Update(sizeMsg)

	gr.merchantMenu = InitializeMerchantMenu(locManager)
	gr.merchantMenu, _ = gr.merchantMenu.Update(sizeMsg)
}

func (gr *GameRender) handleSizeUpdate(msg engine.SizeMsg) {
	gr.screenWidth = msg.Width
	gr.screenHeight = msg.Height

	// Update UI components
	gr.mainMenu, _ = gr.mainMenu.Update(msg)
	gr.classSelection, _ = gr.classSelection.Update(msg)
	gr.settingsMenu, _ = gr.settingsMenu.Update(msg)
	gr.merchantMenu, _ = gr.merchantMenu.Update(msg)
	*gr.hud, _ = gr.hud.Update(msg)

	// Update game space if it exists
	if gr.gameSpace != nil {
		gr.gameSpace.UpdateSize(msg.Width-1, msg.Height-gr.hud.Height()-1)
	}
}

func (gr *GameRender) updateGameSystems() {
	// Update exploration systems
	if gr.gameState.CurrentState == systems.StateExploration {
		// Clean up defeated enemies
		if gr.spawnerSystem != nil {
			gr.spawnerSystem.RemoveDefeatedEnemies()

			// Check if stage is cleared
			if gr.spawnerSystem.IsStageCleared() {
				// Award clearing reward
				if gr.gameInstance != nil && gr.gameInstance.Player != nil && gr.gameInstance.CurrentStage != nil {
					// Add experience or handle stage completion
					// gr.gameInstance.Player.AddExperience(gr.gameInstance.CurrentStage.ClearingReward)
				}
			}
		}
	}

	// Update combat systems
	if gr.gameState.CurrentState == systems.StateCombat {
		if gr.combatSystem != nil && gr.gameInstance != nil && gr.gameInstance.Player != nil {
			gr.combatSystem.Update(gr.gameInstance.Player)
		}
	}
}

// Helper to update HUD stats
func (gr *GameRender) updateHUDStats() {
	if gr.gameInstance == nil || gr.gameInstance.Player == nil {
		return
	}

	player := gr.gameInstance.Player
	worldID := gr.gameInstance.CurrentWorld.WorldID
	stageID := gr.gameInstance.CurrentStage.StageNb

	gr.hud.SetPlayerStats(
		player.Stats.CurrentHP,
		player.Stats.MaxHP,
		player.Stats.Level,
		int(player.Stats.Exp),
		player.Stats.NextLevelExp,
		worldID,
		stageID,
	)

	if gr.gameInstance.CurrentWorld != nil && gr.gameInstance.CurrentStage != nil {
		gr.hud.SetLocation(gr.gameInstance.CurrentWorld.Name, gr.gameInstance.CurrentStage.Name)
	}
}
