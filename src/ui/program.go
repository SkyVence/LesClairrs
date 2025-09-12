package ui

import (
	"fmt"
	"os"
	"time"

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

// TickMsg is a message that is sent on a timer.
type TickMsg struct {
	// The time the tick occurred.
	Time time.Time
}

// Tick is a command that sends a TickMsg after a specified duration.
func Tick(d time.Duration) Cmd {
	return func() Msg {
		time.Sleep(d)
		return TickMsg{Time: time.Now()}
	}
}

// TickNow returns a Tick command that fires immediately
func TickNow() Cmd {
	return func() Msg {
		return TickMsg{Time: time.Now()}
	}
}

type Program struct {
	Model    Model
	renderer renderer
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

func NewProgram(model Model, opts ...ProgramOption) *Program {
	p := &Program{
		Model:    model,
		renderer: newRenderer(os.Stdout),
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

	p.renderer.start()
	defer p.renderer.stop()

	if p.useAltScreen {
		p.renderer.enterAltScreen()
		defer p.renderer.exitAltScreen()
	}

	go readInput(p.msgs)

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
