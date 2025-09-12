package main

import (
	"fmt"
	"log"
	"projectred-rpg.com/ui"
)

type model struct {
	count int
}

var _ ui.Model = &model{}

func newModel() *model {
	return &model{count: 0}
}

func (m *model) Init() ui.Msg {
	return nil
}

func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
	switch msg := msg.(type) {
	case ui.KeyMsg:
		switch msg.Rune {
		case 'q':
			return m, ui.Quit
		case 'j', 'd':
			m.count--
		case 'k', 'u':
			m.count++
		}
	}
	return m, nil
}

func (m *model) View() string {
	return fmt.Sprintf(
		"Count: %d\n\n[j/d] to decrement | [k/u] to increment | [q] to quit",
		m.count,
	)
}


func main() {
	p := ui.NewProgram(newModel(), ui.WithAltScreen())

	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
