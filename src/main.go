package main

import (
	"log"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/game"
)

// main initializes and runs the ProjectRed RPG game engine
func main() {
	g := game.GameModel()

	p := engine.NewProgram(engine.Wrap(g), engine.WithAltScreen())
	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
