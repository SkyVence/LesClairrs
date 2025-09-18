package game

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/loaders"
	"projectred-rpg.com/game/systems"
	"projectred-rpg.com/game/types"
	"projectred-rpg.com/ui"
)

type GameRender struct {
	// Game Systems
	gameInstance  *Game
	gameSpace     *GameRenderer
	gameState     *systems.GameState
	movement      *systems.MovementSystem
	combatSystem  *systems.CombatSystem
	spawnerSystem *systems.SpawnerSystem

	// UI Components
	mainMenu       ui.Menu
	classSelection ui.ClassMenu
	hud            *ui.HUD
	settingsMenu   ui.SettingsMenu
	merchantMenu   ui.MerchantMenu

	// Screen/Renderer Settings
	screenWidth  int
	screenHeight int

	// Input Handling
	inputBuffer []engine.KeyMsg

	// Game Data
	// Add game time
	// Add lang settings
	currentMap *types.TileMap
}

func InitMainMenu(locManager *engine.LocalizationManager) ui.Menu {
	menuOptions := []ui.MenuOption{
		{Label: locManager.Text("ui.menu.start"), Value: "start"},
		{Label: locManager.Text("ui.menu.settings"), Value: "settings"},
		{Label: locManager.Text("ui.menu.quit"), Value: "quit"},
	}
	menu, err := ui.NewMenuWithArtFromFile(locManager.Text("ui.menu.mainmenu"), menuOptions, "assets/logo.txt")
	if err != nil {
		menu = ui.NewMenu("ui.menu.mainmenu", menuOptions)
	}
	return menu
}

func InitializeClassSelection(locManager *engine.LocalizationManager, classes []types.Class) ui.ClassMenu {
	menuOptions := []ui.ClassMenuOption{}
	for _, class := range classes {
		label := class.Name
		if locManager != nil {
			translatedLabel := locManager.Text(class.Name)
			if translatedLabel != "" && translatedLabel != class.Name {
				label = translatedLabel
			}
		}
		menuOptions = append(menuOptions, ui.ClassMenuOption{
			Label: label,
			Value: class.Name,
			Desc:  class.Description,
			Stats: ui.ClassStats{
				MaxHP:    class.MaxHP,
				Force:    class.Force,
				Speed:    class.Speed,
				Defense:  class.Defense,
				Accuracy: class.Accuracy,
			},
		})
	}
	menu := ui.NewClassMenu(locManager.Text("ui.class.menu.name"), menuOptions, locManager)
	return menu
}

func InitializeSettingsSelection(locManager *engine.LocalizationManager, languageOptions []string) ui.SettingsMenu {
	menuOptions := []ui.SettingsMenuOption{}
	for _, lang := range languageOptions {
		menuOptions = append(menuOptions, ui.SettingsMenuOption{
			Label: lang, // Keep the raw language code for localization lookup
			Value: lang,
		})
	}
	menu := ui.NewSettingsMenu(locManager.Text("ui.settings.menu.title"), menuOptions, locManager)
	return menu
}

func InitializeMerchantMenu(locManager *engine.LocalizationManager) ui.MerchantMenu {
	menuOptions := []ui.MerchantMenuOption{
		// === WEAPONS SECTION ===
		{
			Item: types.Item{
				Name:        "═══ WEAPONS ═══",
				Description: "",
				Type:        types.Weapon,
			},
			Price: 0, 
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.katana.name",
				Description: "ui.weapons.katana.description",
				Type:        types.Weapon,
			},
			Price: 150,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.faux neuralink.name",
				Description: "ui.weapons.faux neuralink.description",
				Type:        types.Weapon,
			},
			Price: 200,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.arc synaptique.name",
				Description: "ui.weapons.arc synaptique.description",
				Type:        types.Weapon,
			},
			Price: 175,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.sniper.name",
				Description: "ui.weapons.sniper.description",
				Type:        types.Weapon,
			},
			Price: 250,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.neon reaver.name",
				Description: "ui.weapons.neon reaver.description",
				Type:        types.Weapon,
			},
			Price: 300,
		},
		{
			Item: types.Item{
				Name:        "═══ CONSUMABLE ═══",
				Description: "",
				Type:        types.Consumable,
			},
			Price: 0,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.small_medkit.name",
				Description: "ui.consumable.small_medkit.description",
				Type:        types.Consumable,
			},
			Price: 25,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.large_medkit.name",
				Description: "ui.consumable.large_medkit.description",
				Type:        types.Consumable,
			},
			Price: 50,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.money.name",
				Description: "ui.consumable.money.description",
				Type:        types.Consumable,
			},
			Price: 75,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.serum.name",
				Description: "ui.consumable.serum.description",
				Type:        types.Consumable,
			},
			Price: 100,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.flash.name",
				Description: "ui.consumable.flash.description",
				Type:        types.Consumable,
			},
			Price: 80,
		},
	}

	merchantName := locManager.Text("game.merchants.weapon")
	menu := ui.NewMerchantMenu(merchantName, menuOptions, locManager)
	return menu
}

func initializeGameInstance() *Game {
	// Define default player class - could be moved to config
	defaultClass := types.Class{
		Name:        "null",
		MaxHP:       0,
		Force:       0,
		Speed:       0,
		Defense:     0,
		Accuracy:    0,
		Description: "null",
	}

	return NewGameInstance(defaultClass)
}

func GameModel() *GameRender {
	// Initialize language settings
	locManager := engine.GetLocalizationManager()
	locManager.SetLanguage("fr")

	// Initialize UI Components
	menu := InitMainMenu(locManager)
	hud := ui.NewHud()

	// Init class Select
	classes := config.GetDefaultClasses()
	classSelection := InitializeClassSelection(locManager, classes)
	// Init settings menu
	supportedLanguages, err := locManager.GetSupportedLanguages()
	if err != nil {
		supportedLanguages = []string{"fr"} // Fallback to French
	}
	settingsMenu := InitializeSettingsSelection(locManager, supportedLanguages)
	merchantMenu := InitializeMerchantMenu(locManager)

	// Initialize Game Systems
	gameInstance := initializeGameInstance()
	gameState := systems.NewGameState(systems.StateMainMenu)
	movement := systems.NewMovementSystem()
	spawner := systems.NewSpawnerSystem()
	combatSystem := systems.NewCombatSystem(types.Idle, locManager, spawner)

	return &GameRender{
		gameInstance:  gameInstance,
		gameState:     gameState,
		movement:      movement,
		combatSystem:  combatSystem,
		spawnerSystem: spawner,

		mainMenu:       menu,
		hud:            hud,
		settingsMenu:   settingsMenu,
		classSelection: classSelection,
		merchantMenu:   merchantMenu,

		screenWidth:  80,
		screenHeight: 24,
		inputBuffer:  make([]engine.KeyMsg, 0, 10),
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

func (gr *GameRender) renderGameView() string {
	if gr.gameSpace == nil {
		gr.gameSpace = NewGameRenderer(gr.screenWidth-1, gr.screenHeight-gr.hud.Height()-1)
	}

	// Update HUD with current player stats
	gr.updateHUDStats()

	// Load and set map for current world/stage if available
	if gr.gameInstance != nil && gr.gameInstance.CurrentWorld != nil && gr.gameInstance.CurrentStage != nil {
		tm := loaders.LoadStageMap(gr.gameInstance.CurrentWorld.WorldID, gr.gameInstance.CurrentStage.StageNb)
		gr.currentMap = tm
		gr.gameSpace.SetMap(tm)

		gr.spawnerSystem.LoadStage(gr.gameInstance.CurrentStage)
		// Ensure player spawn is valid for the loaded map
		if gr.gameInstance.Player != nil {
			gr.movement.EnsureValidSpawn(gr.gameInstance.Player, gr.currentMap)
		}
	}

	// Render game world
	if gr.spawnerSystem != nil {
		activeEnemies := gr.spawnerSystem.GetActiveEnemies()
		gr.gameSpace.SetEnemies(activeEnemies)
	}

	gameContent := gr.gameSpace.RenderGameWorld(gr.gameInstance.Player)

	return gr.hud.RenderWithContent(gameContent)
}

func (gr *GameRender) Update(msg engine.Msg) (engine.Model, engine.Cmd) {

	if gr.gameState.CurrentState == systems.StateExploration {
		gr.updateGameSystems()
	}

	switch msg := msg.(type) {
	case engine.SizeMsg:
		gr.handleSizeUpdate(msg)

	case engine.KeyMsg:
		return gr.handleKeyInput(msg)

	default:
		// Handle other message types
	}

	return gr, nil
}

func (gr *GameRender) updateGameSystems() {
	if gr.gameState.CurrentState != systems.StateExploration {
		return
	}

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

func (gr *GameRender) handleKeyInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	currentState := gr.gameState.CurrentState

	switch currentState {
	case systems.StateMainMenu:
		return gr.handleMainMenuInput(msg)

	case systems.StateClassSelection:
		return gr.handleClassSelectionInput(msg)

	case systems.StateSettings:
		return gr.handleSettingsSelectionInput(msg)

	case systems.StateMerchant:
		return gr.handleMerchantInput(msg)

	case systems.StateExploration:
		return gr.handleGameInput(msg)

	default:
		return gr, nil
	}
}

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

func (gr *GameRender) handleGameInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '↑', '↓', '←', '→':
		if gr.gameState.CurrentState == systems.StateExploration {
			// Use movement system with map-based collision and bounds
			_ = gr.movement.MovePlayer(gr.gameInstance.Player, msg.Rune, gr.currentMap)

			if gr.combatSystem.TryEngageCombat(gr.gameInstance.Player) {
				gr.gameState.ChangeState(systems.StateCombat)
			}
		}
	case 'm':
		gr.gameState.ChangeState(systems.StateMerchant)
		gr.merchantMenu, _ = gr.merchantMenu.Update(engine.SizeMsg{Width: gr.screenWidth, Height: gr.screenHeight})
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
			// Find the class by name
			classes := config.GetDefaultClasses()
			for _, class := range classes {
				if class.Name == selected.Value {
					// Initialize game with selected class
					gr.gameInstance = NewGameInstance(class)
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

func (m *GameRender) Init() engine.Msg {
	return nil
}

func (gr *GameRender) View() string {
	if gr.gameState == nil {
		return "Error: Game state is nil"
	}
	currentState := gr.gameState.CurrentState

	switch currentState {
	case systems.StateMainMenu:
		return gr.mainMenu.View()
	case systems.StateClassSelection:
		return gr.classSelection.View()
	case systems.StateExploration:
		return gr.renderGameView()
	case systems.StateSettings:
		return gr.settingsMenu.View()
	case systems.StateMerchant:
		return gr.merchantMenu.View()
	default:
		return "Unknown State"
	}
}
