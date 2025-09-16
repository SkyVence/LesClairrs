package ui

import (
	"time"

	"projectred-rpg.com/engine"
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

func (s Spinner) Init() engine.Cmd {
	return engine.TickNow()
}

func (s Spinner) Update(msg engine.Msg) (Spinner, engine.Cmd) {
	switch msg.(type) {
	case engine.TickMsg:
		s.frame = (s.frame + 1) % len(spinnerFrames)
		return s, engine.Tick(s.speed)
	}
	return s, nil
}

func (s Spinner) View() string {
	return spinnerFrames[s.frame]
}
