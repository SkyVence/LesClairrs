package main

import (
	"log"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/game"
)

func main() {
	gameModel := game.NewGame()

	p := engine.NewProgram(engine.Wrap(gameModel), engine.WithAltScreen())

	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
