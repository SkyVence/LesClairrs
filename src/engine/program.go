package engine

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Program struct {
	Model    Model
	renderer Renderer
	msgs     chan Msg

	useAltScreen bool

	quit bool
}
type ProgramOption func(*Program)

func WithAltScreen() ProgramOption {
	return func(p *Program) {
		p.useAltScreen = true
	}
}

func (p *Program) GetSize() (int, int) {
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		return 80, 24
	}

	width, height, err := term.GetSize(fd)

	if err != nil {
		return 80, 24
	}

	// Sanity check the values
	if width <= 0 || height <= 0 {
		return 80, 24
	}

	return width, height
}

func NewProgram(model Model, opts ...ProgramOption) *Program {
	p := &Program{
		Model:    model,
		renderer: NewRenderer(os.Stdout),
		msgs:     make(chan Msg),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Program) Run() error {

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	p.renderer.Start()
	defer p.renderer.Stop()

	if p.useAltScreen {
		p.renderer.EnterAltScreen()
		defer p.renderer.ExitAltScreen()
	}

	p.renderer.HideCursor()
	go ReadInput(p.msgs)

	// Process the initial message from the model's Init() method.
	var cmd Cmd
	if initialMsg := p.Model.Init(); initialMsg != nil {
		p.Model, cmd = p.Model.Update(initialMsg)
	} else {
		// If Init returns nil, still call Update with nil to trigger initialization
		p.Model, cmd = p.Model.Update(nil)
	}

	// Execute the initial command if it exists
	if cmd != nil {
		go func() {
			p.msgs <- cmd()
		}()
	}

	width, height := p.GetSize()

	p.Model, cmd = p.Model.Update(SizeMsg{Width: width, Height: height})
	if cmd != nil {
		go func() {
			p.msgs <- cmd()
		}()
	}

	for !p.quit {
		// Get the current view from the model.
		view := p.Model.View()

		// Render the view.
		p.renderer.Write(view)

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
