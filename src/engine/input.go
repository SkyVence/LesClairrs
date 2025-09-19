// input.go
package engine

import (
	"os"
)

// ReadInput reads from stdin and sends KeyMsg/QuitMsg to the provided channel
// Handles escape sequences for arrow keys and Ctrl+C termination
func ReadInput(msgs chan<- Msg) {
	buf := make([]byte, 1024)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		data := buf[:n]

		if len(data) >= 3 && data[0] == 0x1b && data[1] == '[' {
			switch data[2] {
			case 'A':
				msgs <- KeyMsg{Rune: '↑'}
				continue
			case 'B':
				msgs <- KeyMsg{Rune: '↓'}
				continue
			case 'C':
				msgs <- KeyMsg{Rune: '→'}
				continue
			case 'D':
				msgs <- KeyMsg{Rune: '←'}
				continue
			}
		}

		switch data[0] {
		case 3:
			msgs <- QuitMsg{}
			return
		case 0x1b:
			if len(data) < 3 {
				continue
			}
			fallthrough
		default:
			if len(data) == 1 {
				msgs <- KeyMsg{Rune: rune(data[0])}
			}
		}
	}
}
