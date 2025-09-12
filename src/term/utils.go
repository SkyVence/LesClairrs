package term

import "fmt"

const (
	Esc = "\x1b"
)

func Clear() {
	fmt.Print(Esc + "[2J")
	fmt.Print(Esc + "[H")
}

func Move(row, col int) {
	fmt.Printf("%s[%d;%dH", Esc, row, col)
}

func HideCursor() { fmt.Print(Esc + "[?25l") }
func ShowCursor() { fmt.Print(Esc + "[?25h") }

func Color(code string) string { return Esc + "[" + code + "m" }
func Reset() string            { return Color("0") }

var (
	Red     = Color("31")
	Green   = Color("32")
	Yellow  = Color("33")
	Blue    = Color("34")
	Magenta = Color("35")
	Cyan    = Color("36")
	White   = Color("37")
	Bold    = Color("1")
)
