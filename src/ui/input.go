// input.go
package ui

import (
	"os"
)

// readInput reads from stdin and sends messages to the provided channel.
// This function runs in a separate goroutine and blocks until input is received.
func readInput(msgs chan<- Msg) {
	// A small buffer to read raw bytes from the terminal.
	buf := make([]byte, 1024)

	for {
		// Read from the standard input.
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		// Get the raw bytes that were read.
		data := buf[:n]

		// The first byte is the most important.
		switch data[0] {
		// Handle Ctrl+C, which sends an End-of-Text character.
		case 3:
			msgs <- QuitMsg{}
			return // Exit the input loop.
		default:
			// For this lightweight version, we'll assume any other single byte
			// is a valid rune. A more advanced parser would handle multi-byte
			// UTF-8 characters and ANSI escape sequences (like arrow keys).
			if len(data) == 1 {
				msgs <- KeyMsg{Rune: rune(data[0])}
			}
		}
	}
}