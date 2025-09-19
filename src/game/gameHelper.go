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

// updateGameSystems handles state-specific system updates for exploration and combat
func (gr *GameRender) updateGameSystems() {
	if gr.gameState.CurrentState == systems.StateExploration {
		if gr.spawnerSystem != nil {
			gr.spawnerSystem.RemoveDefeatedEnemies()

			if gr.spawnerSystem.IsStageCleared() {
				if gr.gameInstance != nil && gr.gameInstance.Player != nil && gr.gameInstance.CurrentStage != nil {

				}
			}
		}
	}

	if gr.gameState.CurrentState == systems.StateCombat {
		if gr.combatSystem != nil && gr.gameInstance != nil && gr.gameInstance.Player != nil {
			gr.combatSystem.Update(gr.gameInstance.Player)

			if gr.combatSystem.IsReadyToExit() {
				if gr.gameInstance.Player.Stats.CurrentHP <= 0 {
					gr.handlePlayerDefeat()
				}

				gr.gameState.ChangeState(systems.StateExploration)

				if gr.gameSpace != nil && gr.spawnerSystem != nil {
					gr.spawnerSystem.RemoveDefeatedEnemies()
					activeEnemies := gr.spawnerSystem.GetActiveEnemies()
					gr.gameSpace.ForceRefreshEnemies(activeEnemies)
				}
			}
		}
	}
}

// updateHUDStats refreshes HUD with current player stats and location info
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

// handlePlayerDefeat respawns player at 25% health in a safe position away from enemies
func (gr *GameRender) handlePlayerDefeat() {
	if gr.gameInstance == nil || gr.gameInstance.Player == nil {
		return
	}

	player := gr.gameInstance.Player

	player.Stats.CurrentHP = player.Stats.MaxHP / 4
	if player.Stats.CurrentHP < 1 {
		player.Stats.CurrentHP = 1
	}

	safePositions := []struct{ x, y int }{
		{2, 2}, {3, 2}, {2, 3}, {3, 3}, {1, 1}, {4, 4},
	}

	for _, pos := range safePositions {
		isSafe := true
		if gr.spawnerSystem != nil {
			for _, enemy := range gr.spawnerSystem.GetActiveEnemies() {
				enemyPos := enemy.GetPosition()
				if abs(enemyPos.X-pos.x) <= 3 && abs(enemyPos.Y-pos.y) <= 3 {
					isSafe = false
					break
				}
			}
		}

		if isSafe {
			player.Pos.X = pos.x
			player.Pos.Y = pos.y
			break
		}
	}
}

// abs returns absolute value of integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
