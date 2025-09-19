package engine

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/ansi"
)

type Renderer interface {
	// Start the renderer
	Start()
	// Stop the renderer
	Stop()
	// Kill the renderer
	Kill()
	// Write to the viewport
	Write(string)
	// Clear the screen
	ClearScreen()
	// Full repaint the screen
	Repaint()
	// Show cursor
	ShowCursor()
	// Hide cursor
	HideCursor()
	// Set window title
	SetWindowTitle(string)
	// Whether or not the alternate screen buffer is enabled.
	AltScreen() bool
	// Enable the alternate screen buffer.
	EnterAltScreen()
	// Disable the alternate screen buffer.
	ExitAltScreen()
	// Position cursor at specific coordinates (0-based)
	SetCursor(x, y int)
	// Get current terminal dimensions
	GetSize() (width int, height int)
}

type StandardRenderer struct {
	mtx *sync.Mutex
	out io.Writer

	buf                bytes.Buffer
	queuedMessageLines []string
	frameRate          time.Duration
	ticker             *time.Ticker
	done               chan struct{}
	lastRender         string
	lastRenderedLines  []string
	linesRendered      int
	altLinesRendered   int
	once               sync.Once

	cursorHidden bool

	altScreenActive bool

	width  int
	height int

	ignoreLines map[int]struct{}
}

// NewRenderer creates a StandardRenderer with default 24fps frameRate
func NewRenderer(out io.Writer) Renderer {
	r := &StandardRenderer{
		out:                out,
		mtx:                &sync.Mutex{},
		done:               make(chan struct{}),
		frameRate:          time.Second / time.Duration(24),
		queuedMessageLines: []string{},
	}
	return r
}

// Start initializes the renderer timer and begins the background rendering loop
func (r *StandardRenderer) Start() {
	if r.ticker == nil {
		r.ticker = time.NewTicker(r.frameRate)
	} else {
		r.ticker.Reset(r.frameRate)
	}

	r.once = sync.Once{}

	go r.listen()
}

// Stop gracefully shuts down the renderer, flushes output, and shows cursor
func (r *StandardRenderer) Stop() {
	r.once.Do(func() {
		r.done <- struct{}{}
	})
	r.flush()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.ShowCursor)
	r.execute(ansi.EraseEntireScreen)
}

// execute writes ANSI escape sequence to output stream
func (r *StandardRenderer) execute(seq string) {
	_, _ = io.WriteString(r.out, seq)
}

// Kill forcefully stops the renderer and clears the current line
func (r *StandardRenderer) Kill() {
	r.once.Do(func() {
		r.done <- struct{}{}
	})

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.EraseEntireLine)
	r.execute("\r")
}

// listen runs the main render loop, waiting for timer ticks or shutdown signal
func (r *StandardRenderer) listen() {
	for {
		select {
		case <-r.done:
			r.ticker.Stop()
			return
		case <-r.ticker.C:
			r.flush()
		}
	}
}

// flush outputs buffered content to terminal with line-by-line optimization
func (r *StandardRenderer) flush() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.buf.Len() == 0 || r.buf.String() == r.lastRender {
		return
	}

	buf := &bytes.Buffer{}

	if r.altScreenActive {
		buf.WriteString(ansi.CursorHomePosition)
	} else if r.linesRendered < 1 {
		buf.WriteString(ansi.CursorUp(r.linesRendered - 1))
	}

	newLines := strings.Split(r.buf.String(), "\n")

	if r.height > 0 && len(newLines) > r.height {
		newLines = newLines[len(newLines)-r.height:]
	}

	flushQueuedMessages := len(r.queuedMessageLines) > 0 && !r.altScreenActive

	if flushQueuedMessages {
		for _, line := range r.queuedMessageLines {
			if ansi.StringWidth(line) > r.width {
				line = line + ansi.EraseLineRight
			}
			_, _ = buf.WriteString(line)
			_, _ = buf.WriteString("\r\n")
		}

		r.queuedMessageLines = []string{}
	}

	for i := 0; i < len(newLines); i++ {
		canSkip := flushQueuedMessages &&
			len(r.lastRenderedLines) > i && r.lastRenderedLines[i] == newLines[i]

		if _, ignore := r.ignoreLines[i]; ignore || canSkip {
			if i < len(newLines)-1 {
				buf.WriteByte('\n')
			}
			continue
		}

		if i == 0 && r.lastRender == "" {
			buf.WriteByte('\r')
		}

		line := newLines[i]

		if r.width > 0 {
			line = ansi.Truncate(line, r.width, "")
		}

		if ansi.StringWidth(line) > r.width {
			line = line + ansi.EraseLineRight
		}

		_, _ = buf.WriteString(line)

		if i < len(newLines)-1 {
			_, _ = buf.WriteString("\r\n")
		}
	}

	if r.lastLinesRendered() > len(newLines) {
		buf.WriteString(ansi.EraseScreenBelow)
	}
	if r.altScreenActive {
		r.altLinesRendered = len(newLines)
	} else {
		r.linesRendered = len(newLines)
	}

	if r.altScreenActive {
		buf.WriteString(ansi.CursorPosition(0, len(newLines)))
	} else {
		buf.WriteByte('\r')
	}

	_, _ = r.out.Write(buf.Bytes())
	r.lastRender = r.buf.String()

	r.lastRenderedLines = newLines
	r.buf.Reset()
}

// lastLinesRendered returns appropriate line count based on screen mode
func (r *StandardRenderer) lastLinesRendered() int {
	if r.altScreenActive {
		return r.altLinesRendered
	}
	return r.linesRendered
}

// ClearScreen erases entire screen and forces repaint
func (r *StandardRenderer) ClearScreen() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.EraseEntireScreen)
	r.execute(ansi.CursorHomePosition)

	r.Repaint()
}

// Repaint resets render state to force full redraw on next flush
func (r *StandardRenderer) Repaint() {
	r.lastRender = ""
	r.lastRenderedLines = nil
}

// AltScreen returns whether alternate screen buffer is currently active
func (r *StandardRenderer) AltScreen() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.altScreenActive
}

// EnterAltScreen switches to alternate screen buffer for full-screen mode
func (r *StandardRenderer) EnterAltScreen() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.altScreenActive {
		return
	}

	r.altScreenActive = true
	r.execute(ansi.SetAltScreenSaveCursorMode)

	r.execute(ansi.EraseEntireScreen)
	r.execute(ansi.CursorHomePosition)

	if r.cursorHidden {
		r.execute(ansi.HideCursor)
	} else {
		r.execute(ansi.ShowCursor)
	}

	r.altLinesRendered = 0

	r.Repaint()
}

// ExitAltScreen restores normal screen buffer and cursor position
func (r *StandardRenderer) ExitAltScreen() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if !r.altScreenActive {
		return
	}

	r.altScreenActive = false
	r.execute(ansi.ResetAltScreenSaveCursorMode)

	if r.cursorHidden {
		r.execute(ansi.HideCursor)
	} else {
		r.execute(ansi.ShowCursor)
	}

	r.Repaint()
}

// ShowCursor makes the terminal cursor visible
func (r *StandardRenderer) ShowCursor() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.cursorHidden = false
	r.execute(ansi.ShowCursor)
}

// HideCursor makes the terminal cursor invisible
func (r *StandardRenderer) HideCursor() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.cursorHidden = true
	r.execute(ansi.HideCursor)
}

// SetWindowTitle sets the terminal window title
func (r *StandardRenderer) SetWindowTitle(title string) {
	r.execute(ansi.SetWindowTitle(title))
}

// Write stores content in render buffer, replacing any previous content
func (r *StandardRenderer) Write(s string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.buf.Reset()

	if s == "" {
		s = " "
	}

	_, _ = r.buf.WriteString(s)
}

// SetCursor positions cursor at specific coordinates in alternate screen mode
func (r *StandardRenderer) SetCursor(x, y int) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.altScreenActive {
		r.execute(ansi.CursorPosition(y, x))
	}
}

// GetSize returns current terminal width and height
func (r *StandardRenderer) GetSize() (int, int) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.width, r.height
}
