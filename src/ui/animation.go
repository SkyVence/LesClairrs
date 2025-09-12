package ui

import (
	"errors"
	"os"
	"strings"
	"time"
)

func LoadAnimationFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("Animation file is empty")
	}

	rawFrames := strings.Split(string(content), "---")

	var cleanedFrames []string
	for _, frame := range rawFrames {
		trimmedFrame := strings.TrimSpace(frame)
		if trimmedFrame != "" {
			cleanedFrames = append(cleanedFrames, trimmedFrame)
		}
	}

	if len(cleanedFrames) == 0 {
		return nil, errors.New("No valid frames found in animation file")
	}

	return cleanedFrames, nil
}

type Animation struct {
	Frames []string
	frame  int
	speed  time.Duration
}

func NewAnimation(frames []string) Animation {
	return Animation{
		Frames: frames,
		speed:  200 * time.Millisecond,
	}
}

func (a Animation) Init() Cmd {
	return Tick(a.speed)
}

func (a Animation) Update(msg Msg) (Animation, Cmd) {
	switch msg.(type) {
	case TickMsg:
		if len(a.Frames) > 0 {
			a.frame = (a.frame + 1) % len(a.Frames)
		}
		return a, Tick(a.speed)
	}
	return a, nil
}

func (a Animation) View() string {
	if len(a.Frames) == 0 {
		return "Animation has no frames."
	}
	return a.Frames[a.frame]
}
