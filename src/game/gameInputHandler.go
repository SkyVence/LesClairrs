package game

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/systems"
	"projectred-rpg.com/game/types"
)

func (gr *GameRender) handleGameInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '↑', '↓', '←', '→':
		if gr.gameState.CurrentState == systems.StateExploration {
			_ = gr.movement.MovePlayer(gr.gameInstance.Player, msg.Rune, gr.currentMap)

			if gr.combatSystem.TryEngageCombat(gr.gameInstance.Player) {
				gr.gameState.ChangeState(systems.StateCombat)
			}
		} else if gr.gameState.CurrentState == systems.StateCombat {
			gr.handleCombatInput(msg)
		}
	case 'm':
		gr.gameState.ChangeState(systems.StateMerchant)
		return gr, nil
	case 'p':
		if gr.gameState.CurrentState == systems.StateExploration {
			gr.transitionToNextLevel()
		}
		return gr, nil
	}

	return gr, nil
}

func (gr *GameRender) handleMerchantInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '\r', '\n', ' ':
		return gr, nil
	case 'q':
		gr.gameState.ChangeState(systems.StateExploration)
		return gr, nil
	default:
		gr.merchantMenu, _ = gr.merchantMenu.Update(msg)
	}
	return gr, nil
}

func (gr *GameRender) handleClassSelectionInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '\r', '\n', ' ': // Enter key
		selected := gr.classSelection.GetSelected()
		if selected.Value != "" {
			classes := config.GetDefaultClasses()
			for _, class := range classes {
				if class.Name == selected.Value {
					// Get current language
					currentLang := engine.GetLocalizationManager().GetCurrentLanguage()

					// Initialize game with selected class and language
					gr.gameInstance = NewGameInstance(class, currentLang) // CORRIGÉ

					gr.gameInstance.LoadStage(1, 1)

					gr.gameState.ChangeState(systems.StateExploration)
					return gr, nil
				}
			}
		}

	case 'q':
		gr.gameState.ChangeState(systems.StateMainMenu)
		// Ensure main menu has current dimensions when returning
		gr.mainMenu, _ = gr.mainMenu.Update(engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight})
		return gr, nil

	default:
		// Pass input to menu for navigation
		gr.classSelection, _ = gr.classSelection.Update(msg)
	}
	return gr, nil
}

func (gr *GameRender) handleMainMenuInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '\r', '\n', ' ': // Enter key
		selected := gr.mainMenu.GetSelected()
		switch selected.Value {
		case "start":
			// Transition to class selection
			gr.gameState.ChangeState(systems.StateClassSelection)
			gr.classSelection, _ = gr.classSelection.Update(engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight})

			return gr, nil

		case "settings":
			gr.gameState.ChangeState(systems.StateSettings)
			gr.settingsMenu, _ = gr.settingsMenu.Update(engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight})
			return gr, nil

		case "quit":
			return gr, engine.Quit
		}

	case 'q':
		return gr, engine.Quit

	default:
		// Pass input to menu for navigation
		gr.mainMenu, _ = gr.mainMenu.Update(msg)
	}

	return gr, nil
}

func (gr *GameRender) handleSettingsSelectionInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '\r', '\n', ' ':
		selected := gr.settingsMenu.GetSelected()
		if selected.Value != "" {
			err := engine.GetLocalizationManager().SetLanguage(selected.Value)
			if err == nil {
				gr.refreshMenusAfterLanguageChange()
			}
			gr.gameState.ChangeState(systems.StateMainMenu)
			return gr, nil
		}
	case 'q':
		gr.gameState.ChangeState(systems.StateMainMenu)
		gr.mainMenu, _ = gr.mainMenu.Update(engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight})
		return gr, nil
	default:
		// Pass input to menu for navigation
		gr.settingsMenu, _ = gr.settingsMenu.Update(msg)
	}
	return gr, nil
}

// handleCombatInput handles input during combat state
func (gr *GameRender) handleCombatInput(msg engine.KeyMsg) {
	if gr.combatSystem.CurrentCombatState != types.PlayerTurn {
		return // Only handle input during player turn
	}

	combatUI := gr.combatSystem.GetCombatUI()

	switch msg.Rune {
	case '↑':
		combatUI.SelectPrevAction()
	case '↓':
		combatUI.SelectNextAction()
	case '\r', '\n', ' ': // Enter or Space - confirm action
		action := combatUI.GetSelectedAction()
		if action != "" {
			success := gr.combatSystem.ProcessPlayerAction(action, gr.gameInstance.Player)
			if action == "Run" && success {
				// Successfully ran away, return to exploration
				gr.gameState.ChangeState(systems.StateExploration)
			}
		}
	case 'h', 'H': // Toggle history
		combatUI.ToggleHistory()
	case 'q', 'Q': // Force quit combat (emergency exit)
		gr.combatSystem.ExitCombat()
		gr.gameState.ChangeState(systems.StateExploration)
	}
}
