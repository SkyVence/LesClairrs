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
	locManager    *engine.LocalizationManager

	// UI Components
	hud            *ui.HUD
	mainMenu       ui.Menu
	classSelection ui.ClassMenu
	settingsMenu   ui.SettingsMenu
	merchantMenu   ui.MerchantMenu
	combatHud      *ui.CombatHUD

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
	gameInstance := initializeGameInstance() // Pass current language
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
		locManager:    locManager,

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
	switch msg := msg.(type) {
	case engine.SizeMsg:
		gr.handleSizeUpdate(msg)

	case engine.KeyMsg:
		return gr.handleKeyInput(msg)

	case engine.TickMsg:
		// Handle level intro tick updates
		if gr.gameInstance != nil && gr.gameInstance.IsShowingIntro() {
			var cmd engine.Cmd
			gr.gameInstance.LevelIntro, cmd = gr.gameInstance.LevelIntro.Update(msg)
			return gr, cmd
		}

	default:
		// Handle other message types
	}

	return gr, nil
}

func (gr *GameRender) handleKeyInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	currentState := gr.gameState.CurrentState

	// Handle level intro separately
	//if gr.gameInstance != nil && gr.gameInstance.IsShowingIntro() {
	//	return gr.handleLevelIntroInput(msg)
	//}

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

func (m *GameRender) Init() engine.Msg {
	return nil
}

func (gr *GameRender) View() string {
	if gr.gameState == nil {
		return "Error: Game state is nil"
	}

	// If showing level intro, render it over everything
	if gr.gameInstance != nil && gr.gameInstance.IsShowingIntro() {
		return gr.gameInstance.LevelIntro.Render()
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
