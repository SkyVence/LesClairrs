package game

import (
	"projectred-rpg.com/config"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/systems"
	"projectred-rpg.com/game/types"
	"projectred-rpg.com/ui"
)

type GameRender struct {
	// Game Systems
	gameInstance *Game
	gameSpace    *GameRenderer
	gameState    *systems.GameState

	// UI Components
	mainMenu       ui.Menu
	classSelection ui.ClassMenu
	hud            *ui.HUD

	// Screen/Renderer Settings
	screenWidth  int
	screenHeight int

	// Input Handling
	inputBuffer []engine.KeyMsg

	// Game Data
	// Add game time
	// Add lang settings
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

func initializeGameInstance() *Game {
	// Define default player class - could be moved to config
	defaultClass := types.Class{
		Name:        "Cyber-Samurai",
		MaxHP:       100,
		Force:       1,
		Speed:       12,
		Defense:     8,
		Accuracy:    15,
		Description: "A swift and deadly warrior, excelling in close combat and agility.",
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
	gameInstance := initializeGameInstance()
	gameState := systems.NewGameState(systems.StateMainMenu)
	return &GameRender{
		gameInstance: gameInstance,
		gameState:    gameState,

		mainMenu:       menu,
		hud:            hud,
		classSelection: classSelection,

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

	// Render game world
	gameContent := gr.gameSpace.RenderGameWorld(gr.gameInstance.Player)

	return gr.hud.RenderWithContent(gameContent)
}

func (gr *GameRender) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
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

func (gr *GameRender) handleSizeUpdate(msg engine.SizeMsg) {
	gr.screenWidth = msg.Width
	gr.screenHeight = msg.Height

	// Update UI components
	gr.mainMenu, _ = gr.mainMenu.Update(msg)
	gr.classSelection, _ = gr.classSelection.Update(msg)
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

	case systems.StateExploration:
		return gr.handleGameInput(msg)

	default:
		return gr, nil
	}
}

func (gr *GameRender) handleGameInput(msg engine.KeyMsg) (engine.Model, engine.Cmd) {
	switch msg.Rune {
	case '↑', '↓', '←', '→':
		if gr.gameState.CurrentState == systems.StateExploration {
			gr.gameInstance.Player.Move(msg.Rune, gr.gameSpace.width, gr.gameSpace.height)
		}
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
		return "Settings Menu - (Not Implemented)"
	default:
		return "Unknown State"
	}
}
