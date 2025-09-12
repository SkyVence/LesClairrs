package loop

import (
	"bufio"
	"os"
	"time"

	"golang.org/x/term"
)

type KeyMsg struct {
	Rune rune
	Seq  string // raw escape sequence for arrows, etc.
}

type TickMsg struct{}

func ReadInput(ch chan<- interface{}) {
	// Put terminal into raw mode
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			close(ch)
			return
		}
		if b == 0x1b { // ESC sequence
			reader.ReadByte() // '['
			reader.ReadByte() // e.g., 'A'
			// Minimal example; you’d parse properly
			ch <- KeyMsg{Seq: "ESC"}
			continue
		}
		ch <- KeyMsg{Rune: rune(b)}
	}
}

func Ticker(d time.Duration, ch chan<- interface{}) {
	t := time.NewTicker(d)
	defer t.Stop()
	for range t.C {
		ch <- TickMsg{}
	}
}
