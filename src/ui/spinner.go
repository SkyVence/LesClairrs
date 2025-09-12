package ui

import (
	"time"
)

const DefaultSpinnerInterval = 100 * time.Millisecond

var spinnerFrames = []string{"/", "-", "\\", "|"}

type Spinner struct {
	frame int
	speed time.Duration
}

func NewSpinner() Spinner {
	return Spinner{
		speed: DefaultSpinnerInterval,
	}
}

func (s Spinner) Init() Cmd {
	return TickNow()
}

func (s Spinner) Update(msg Msg) (Spinner, Cmd) {
	switch msg.(type) {
	case TickMsg:
		s.frame = (s.frame + 1) % len(spinnerFrames)
		return s, Tick(s.speed)
	}
	return s, nil
}

func (s Spinner) View() string {
	return spinnerFrames[s.frame]
}
