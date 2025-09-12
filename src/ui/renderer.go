package ui

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/ansi"
)

type renderer interface {
	// Start the renderer
	start()
	// Stop the renderer
	stop()
	// Kill the renderer
	kill()
	// Write to the viewport
	write(string)
	// Clear the screen
	clearScreen()
	// Full repaint the screen
	repaint()
	// Show cursor
	showCursor()
	// Hide cursor
	hideCursor()
	// Set window title
	setWindowTitle(string)
	// Whether or not the alternate screen buffer is enabled.
	altScreen() bool
	// Enable the alternate screen buffer.
	enterAltScreen()
	// Disable the alternate screen buffer.
	exitAltScreen()
}

type standardRenderer struct {
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

func newRenderer(out io.Writer) renderer {
	r := &standardRenderer{
		out:                out,
		mtx:                &sync.Mutex{},
		done:               make(chan struct{}),
		frameRate:          time.Second / time.Duration(45),
		queuedMessageLines: []string{},
	}
	return r
}

func (r *standardRenderer) start() {
	if r.ticker == nil {
		r.ticker = time.NewTicker(r.frameRate)
	} else {
		r.ticker.Reset(r.frameRate)	
	}

	r.once = sync.Once{}

	go r.listen()
}

func (r *standardRenderer) stop() {
	r.once.Do(func() {
		r.done <- struct{}{}
	})
	r.flush()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.EraseEntireLine)
	r.execute("\r")
}

func (r *standardRenderer) execute(seq string) {
	_, _ = io.WriteString(r.out, seq)
}

func (r *standardRenderer) kill() {
	r.once.Do(func() {
		r.done <- struct{}{}
	})

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.EraseEntireLine)
	r.execute("\r")
}

func (r *standardRenderer) listen() {
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

func (r *standardRenderer) flush() {
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

func (r *standardRenderer) lastLinesRendered() int {
	if r.altScreenActive {
		return r.altLinesRendered
	}
	return r.linesRendered
}

func (r *standardRenderer) clearScreen() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.execute(ansi.EraseEntireScreen)
	r.execute(ansi.CursorHomePosition)

	r.repaint()
}

func (r *standardRenderer) repaint() {
	r.lastRender = ""
	r.lastRenderedLines = nil
}

func (r *standardRenderer) altScreen() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.altScreenActive
}

func (r *standardRenderer) enterAltScreen() {
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

	r.repaint()
}

func (r *standardRenderer) exitAltScreen() {
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

	r.repaint()
}

func (r *standardRenderer) showCursor() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.cursorHidden = false
	r.execute(ansi.ShowCursor)
}

func (r *standardRenderer) hideCursor() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.cursorHidden = true
	r.execute(ansi.HideCursor)
}

func (r *standardRenderer) setWindowTitle(title string) {
	r.execute(ansi.SetWindowTitle(title))
}

func (r *standardRenderer) write(s string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.buf.Reset()

	if s == "" {
		s = " "
	}

	_, _ = r.buf.WriteString(s)
}
