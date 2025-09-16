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

type SizeMsg struct {
	Width  int
	Height int
}

func (p *Program) GetSize() (int, int) {
	fd := int(os.Stdin.Fd())
	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Getting terminal size for fd=%d\n", fd)

	// Check environment variables that might affect terminal detection
	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: TERM=%s\n", os.Getenv("TERM"))
	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: COLUMNS=%s\n", os.Getenv("COLUMNS"))
	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: LINES=%s\n", os.Getenv("LINES"))

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: stdin (fd=%d) is not a terminal!\n", fd)
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Returning default size: 80x24\n")
		return 80, 24
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: stdin is a valid terminal\n")

	width, height, err := term.GetSize(fd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Error getting terminal size: %v\n", err)
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Error type: %T\n", err)
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Returning default size: 80x24\n")
		return 80, 24
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: Successfully got terminal size: width=%d, height=%d\n", width, height)

	// Sanity check the values
	if width <= 0 || height <= 0 {
		fmt.Fprintf(os.Stderr, "[DEBUG] GetSize: WARNING: Invalid size values (width=%d, height=%d), returning defaults\n", width, height)
		return 80, 24
	}

	return width, height
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
	fmt.Fprintf(os.Stderr, "[DEBUG] Run: Starting program\n")

	fd := int(os.Stdin.Fd())
	fmt.Fprintf(os.Stderr, "[DEBUG] Run: stdin fd=%d\n", fd)
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

	p.renderer.hideCursor()
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

	fmt.Fprintf(os.Stderr, "[DEBUG] Run: About to call GetSize\n")
	width, height := p.GetSize()
	fmt.Fprintf(os.Stderr, "[DEBUG] Run: GetSize returned width=%d, height=%d\n", width, height)

	fmt.Fprintf(os.Stderr, "[DEBUG] Run: Sending SizeMsg to model with width=%d, height=%d\n", width, height)
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
