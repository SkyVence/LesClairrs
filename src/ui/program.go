package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Msg interface{}

type Cmd func() Msg

type Model interface {
	Init() Msg
	Update(msg Msg) (Model, Cmd)
	View() string
}

type KeyMsg struct {
	Rune rune
}

type QuitMsg struct{}

func Quit() Msg {
	return QuitMsg{}
}

type Program struct {
	Model    Model
	renderer renderer
	msgs     chan Msg
	quit     bool
}

func NewProgram(model Model) *Program {
	renderer := newRenderer(os.Stdout)

	return &Program{
		Model:    model,
		renderer: renderer,
		msgs:     make(chan Msg),
	}
}

func (p *Program) Run() error {

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	p.renderer.start()
	defer p.renderer.stop()

	go readInput(p.msgs)
	// Process the initial message from the model's Init() method.
	if initialMsg := p.Model.Init(); initialMsg != nil {
		p.Model, _ = p.Model.Update(initialMsg)
	}

	for !p.quit {
		// Get the current view from the model.
		view := p.Model.View()

		// Render the view.
		p.renderer.write(view)

		// Wait for the next message from any source (e.g., keyboard input).
		msg := <-p.msgs

		// Handle the quit message specifically.
		if _, ok := msg.(QuitMsg); ok {
			p.quit = true
			return nil
		}

		// Process the message by calling the model's Update function.
		var cmd Cmd
		p.Model, cmd = p.Model.Update(msg)

		if cmd != nil {
			go func() {
				p.msgs <- cmd()
			}()
		}
	}

	return nil
}
