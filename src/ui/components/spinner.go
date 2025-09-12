package components

import (
	"time"

	"projectred-rpg.com/ui"
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

func (s Spinner) Init() ui.Cmd {
	return ui.TickNow()
}

func (s Spinner) Update(msg ui.Msg) (Spinner, ui.Cmd) {
	switch msg.(type) {
	case ui.TickMsg:
		s.frame = (s.frame + 1) % len(spinnerFrames)
		return s, ui.Tick(s.speed)
	}
	return s, nil
}

func (s Spinner) View() string {
	return spinnerFrames[s.frame]
}
