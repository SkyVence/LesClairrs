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

	// Update combat UI if it exists
	if gr.combatSystem.GetCombatUI() != nil {
		gr.combatSystem.GetCombatUI().Update(msg)
	}

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

			// Check if combat is ready to exit
			if gr.combatSystem.IsReadyToExit() {
				// Check if player was defeated and handle respawn
				if gr.gameInstance.Player.Stats.CurrentHP <= 0 {
					gr.handlePlayerDefeat()
				}

				// Combat has ended, return to exploration
				gr.gameState.ChangeState(systems.StateExploration)

				// Force a refresh of the game space after state transition
				if gr.gameSpace != nil && gr.spawnerSystem != nil {
					// Ensure defeated enemies are removed and refresh the enemy list
					gr.spawnerSystem.RemoveDefeatedEnemies()
					activeEnemies := gr.spawnerSystem.GetActiveEnemies()
					gr.gameSpace.ForceRefreshEnemies(activeEnemies)
				}
			}
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

// handlePlayerDefeat manages player respawn after defeat
func (gr *GameRender) handlePlayerDefeat() {
	if gr.gameInstance == nil || gr.gameInstance.Player == nil {
		return
	}

	player := gr.gameInstance.Player

	// Restore player to 25% health
	player.Stats.CurrentHP = player.Stats.MaxHP / 4
	if player.Stats.CurrentHP < 1 {
		player.Stats.CurrentHP = 1
	}

	// Move player to a safe spawn position (usually 1,1 or 2,2)
	safePositions := []struct{ x, y int }{
		{2, 2}, {3, 2}, {2, 3}, {3, 3}, {1, 1}, {4, 4},
	}

	// Find a safe position away from enemies
	for _, pos := range safePositions {
		// Check if this position is safe (no enemies nearby)
		isSafe := true
		if gr.spawnerSystem != nil {
			for _, enemy := range gr.spawnerSystem.GetActiveEnemies() {
				enemyPos := enemy.GetPosition()
				// Check if enemy is within 3 cells of this position
				if abs(enemyPos.X-pos.x) <= 3 && abs(enemyPos.Y-pos.y) <= 3 {
					isSafe = false
					break
				}
			}
		}

		if isSafe {
			// Set player position
			player.Pos.X = pos.x
			player.Pos.Y = pos.y
			break
		}
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
