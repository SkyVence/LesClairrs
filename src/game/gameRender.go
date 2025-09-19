package game

import (
	"time"

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

	// USE YOUR ENGINE TO GET THE LANGUAGE
	locManager := engine.GetLocalizationManager()
	currentLang := locManager.GetCurrentLanguage()

	return NewGameInstance(defaultClass, currentLang)
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

	gr.updateGameSystems()

	// Update UI components based on message type
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

		// Keep ticking during combat to process enemy turns
		if gr.gameState.CurrentState == systems.StateCombat {
			return gr, engine.Tick(time.Second / 60) // 60 FPS tick rate
		}

	default:
		// Handle other message types
	}

	return gr, nil
}

func (gr *GameRender) handleKeyInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	currentState := gr.gameState.CurrentState

	// Handle level intro separately - DÉCOMMENTEZ CES LIGNES !
	if gr.gameInstance != nil && gr.gameInstance.IsShowingIntro() {
		return gr.handleLevelIntroInput(msg)
	}

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
	case systems.StateCombat:
		return gr.handleGameInput(msg)

	default:
		return gr, nil
	}
}

func (gr *GameRender) handleLevelIntroInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	if gr.gameInstance == nil || gr.gameInstance.LevelIntro == nil {
		return gr, nil
	}

	// Update the intro system
	var cmd engine.Cmd
	gr.gameInstance.LevelIntro, cmd = gr.gameInstance.LevelIntro.Update(msg)

	return gr, cmd
}

func (m *GameRender) Init() engine.Msg {
	// Initialize the combat UI renderer
	renderer := engine.GetGlobalRenderer()
	if renderer != nil {
		m.SetRenderer(renderer)

		// Ensure combat UI gets proper initial size
		if m.combatSystem.GetCombatUI() != nil {
			sizeMsg := engine.SizeMsg{Width: m.screenWidth, Height: m.screenHeight}
			m.combatSystem.GetCombatUI().Update(sizeMsg)
		}
	}
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
	case systems.StateCombat:
		if gr.combatSystem.GetCombatUI() != nil {
			return gr.combatSystem.GetCombatUI().View()
		}
		return "Combat UI not initialized"
	case systems.StateMerchant:
		return gr.merchantMenu.View()
	default:
		return "Unknown State"
	}
}

// SetRenderer initializes the combat UI with the provided renderer
func (gr *GameRender) SetRenderer(renderer engine.Renderer) {
	gr.combatSystem.SetRenderer(renderer)

	// Set up combat exit callback to refresh gameSpace
	gr.combatSystem.SetExitCallback(func() {
		if gr.gameSpace != nil && gr.spawnerSystem != nil {
			activeEnemies := gr.spawnerSystem.GetActiveEnemies()
			gr.gameSpace.ForceRefreshEnemies(activeEnemies)
		}
	})
}

// À ajouter dans votre fichier GameRender principal
func (gr *GameRender) transitionToNextLevel() {
	if gr.gameInstance == nil || gr.gameInstance.CurrentWorld == nil || gr.gameInstance.CurrentStage == nil {
		return
	}

	currentWorldID := gr.gameInstance.CurrentWorld.WorldID
	currentStageNb := gr.gameInstance.CurrentStage.StageNb

	nextStageNb := currentStageNb + 1

	stageExists := false
	for _, stage := range gr.gameInstance.CurrentWorld.Stages {
		if stage.StageNb == nextStageNb {
			stageExists = true
			break
		}
	}

	if stageExists {
		gr.gameInstance.LoadStage(currentWorldID, nextStageNb)
	} else {
		nextWorldID := currentWorldID + 1

		if world, exists := loaders.GetWorld(nextWorldID); exists && len(world.Stages) > 0 {
			gr.gameInstance.LoadStage(nextWorldID, 1)
		} else {
			return
		}
	}

	// Forcer plusieurs mouvements pour sortir le joueur d'une collision
	for i := 0; i < 5; i++ {
		gr.movement.MovePlayer(gr.gameInstance.Player, '→', gr.currentMap)
		gr.movement.MovePlayer(gr.gameInstance.Player, '↓', gr.currentMap)
	}
}

func (gr *GameRender) fixPlayerPosition() {
	// Essayer quelques positions de base qui sont généralement libres
	safePositions := []struct{ x, y int }{
		{1, 1}, {2, 1}, {1, 2}, {2, 2}, {3, 1}, {1, 3}, {3, 3}, {4, 4}, {5, 5},
	}

	for _, pos := range safePositions {
		// Utiliser le système de mouvement pour positionner le joueur
		// en utilisant une direction factice pour déclencher le positionnement
		originalPos := gr.getPlayerPosition()

		if gr.setPlayerPosition(pos.x, pos.y) {
			// Vérifier si cette position est valide en essayant un mouvement nul
			if gr.movement.MovePlayer(gr.gameInstance.Player, '↑', gr.currentMap) {
				// Si le mouvement est valide, on reste ici
				gr.setPlayerPosition(pos.x, pos.y)
				return
			}
		}

		// Restaurer la position originale si celle-ci ne fonctionne pas
		gr.setPlayerPosition(originalPos.x, originalPos.y)
	}
}

// Fonctions helper (à adapter selon votre structure Player)
func (gr *GameRender) getPlayerPosition() struct{ x, y int } {
	// À adapter selon vos champs Player réels
	return struct{ x, y int }{x: 1, y: 1} // position par défaut
}

func (gr *GameRender) setPlayerPosition(x, y int) bool {
	// À adapter selon vos champs Player réels
	// Retourner true si le positionnement a réussi
	return true
}
