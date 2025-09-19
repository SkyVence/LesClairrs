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

// WithAltScreen enables alternate screen buffer for full-screen applications
func WithAltScreen() ProgramOption {
	return func(p *Program) {
		p.useAltScreen = true
	}
}

// GetSize returns terminal width and height, defaulting to 80x24 for non-terminals
func (p *Program) GetSize() (int, int) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return 80, 24
	}

	width, height, err := term.GetSize(fd)

	if err != nil {
		return 80, 24
	}

	if width <= 0 || height <= 0 {
		return 80, 24
	}

	return width, height
}

// GetRenderer returns the renderer instance for external access
func (p *Program) GetRenderer() Renderer {
	return p.renderer
}

// NewProgram creates a new Program with model and applies provided options
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

// Run starts the program main loop, setting up terminal and handling input/rendering
func (p *Program) Run() error {

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	p.renderer.Start()
	defer p.renderer.Stop()

	SetGlobalRenderer(p.renderer)

	if p.useAltScreen {
		p.renderer.EnterAltScreen()
		defer p.renderer.ExitAltScreen()
	}

	p.renderer.HideCursor()
	go ReadInput(p.msgs)

	var cmd Cmd
	if initialMsg := p.Model.Init(); initialMsg != nil {
		p.Model, cmd = p.Model.Update(initialMsg)
	} else {
		p.Model, cmd = p.Model.Update(nil)
	}

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
		view := p.Model.View()

		p.renderer.Write(view)

		msg := <-p.msgs

		if _, ok := msg.(QuitMsg); ok {
			p.quit = true
			return nil
		}

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
