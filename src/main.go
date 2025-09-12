package main

import (
	"fmt"
	"time"

	"projectred-rpg.com/components"
	"projectred-rpg.com/loop"
	"projectred-rpg.com/term"
	"projectred-rpg.com/ui"
)

func main() {
	term.HideCursor()
	defer term.ShowCursor()

	// Root component: a Box containing a Menu
	menu := &components.Menu{Items: []string{"Start", "Options", "Quit"}}
	root := &components.Box{
		Title:   "Project Red: RPG",
		Style:   components.DefaultBoxStyle,
		Padding: 1,
		Child:   menu,
	}
	root.SetBounds(ui.Rect{X: 2, Y: 2, W: 40, H: 10})
	cmds := root.Init()

	renderer := ui.Renderer{}
	msgCh := make(chan any, 16)

	// Start input and ticker
	go loop.ReadInput(msgCh)
	go loop.Ticker(100*time.Millisecond, msgCh)

	// Run initial commands
	for _, c := range cmds {
		go func(cf ui.Cmd) { msgCh <- cf() }(c)
	}

	// Simple controller translating keys to messages
	go func() {
		for m := range msgCh {
			switch v := m.(type) {
			case loop.KeyMsg:
				switch v.Rune {
				case 'z':
					msgCh <- components.MoveUp{}
				case 's':
					msgCh <- components.MoveDown{}
				case '\r':
					msgCh <- components.Enter{}
				case 'q':
					close(msgCh)
					return
				}
			}

			root.Update(m)
			renderer.Render(root.View())
		}
	}()

	renderer.Render(root.View())

	for range msgCh {
	}
	fmt.Println("\nBye!")
}
