package main

import (
	"log"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/game"
	"projectred-rpg.com/ui"
)

func main() {
	// Create the game instance (game logic lives in the game package)
	g := game.NewGame()

	// Wrap it with the engine adapter so it satisfies ui.Model
	p := ui.NewProgram(engine.Wrap(g), ui.WithAltScreen())
	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
