package game

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/systems"
)

func (gr *GameRender) handleGameInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '↑', '↓', '←', '→':
		if gr.gameState.CurrentState == systems.StateExploration {
			_ = gr.movement.MovePlayer(gr.gameInstance.Player, msg.Rune, gr.currentMap)

			if gr.combatSystem.TryEngageCombat(gr.gameInstance.Player) {
				gr.gameState.ChangeState(systems.StateCombat)
			}
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
