package engine

import "time"

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

type SizeMsg struct {
	Width  int
	Height int
}
