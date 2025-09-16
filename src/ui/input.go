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

		// Handle escape sequences (like arrow keys)
		if len(data) >= 3 && data[0] == 0x1b && data[1] == '[' {
			switch data[2] {
			case 'A':
				msgs <- KeyMsg{Rune: '↑'} // Or use a special rune/code
				continue
			case 'B':
				msgs <- KeyMsg{Rune: '↓'} // Down arrow
				continue
			case 'C':
				msgs <- KeyMsg{Rune: '→'} // Right arrow
				continue
			case 'D':
				msgs <- KeyMsg{Rune: '←'} // Left arrow
				continue
			}
		}

		// Handle single byte inputs
		switch data[0] {
		// Handle Ctrl+C, which sends an End-of-Text character.
		case 3:
			msgs <- QuitMsg{}
			return // Exit the input loop.
		// Skip incomplete escape sequences
		case 0x1b:
			if len(data) < 3 {
				continue // Wait for complete sequence
			}
			// If we get here, it's an unhandled escape sequence
			// Fall through to default case
			fallthrough
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