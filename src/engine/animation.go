package engine

import (
	"errors"
	"os"
	"strings"
	"time"
)

// LoadAnimationFile reads and parses animation frames from file separated by "---"
// Returns slice of frame strings and any file read/parse error
func LoadAnimationFile(filename string) ([]string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("animation file is empty")
	}

	sanitizedContent := strings.ReplaceAll(string(content), "\r", "")
	rawFrames := strings.Split(sanitizedContent, "---")

	var cleanedFrames []string
	for _, frame := range rawFrames {
		trimmedFrame := strings.Trim(frame, "\n\t")
		if trimmedFrame != "" {
			cleanedFrames = append(cleanedFrames, trimmedFrame)
		}
	}

	if len(cleanedFrames) == 0 {
		return nil, errors.New("no valid frames found in animation file")
	}

	return cleanedFrames, nil
}

type Animation struct {
	Frames []string
	frame  int
	speed  time.Duration
}

// NewAnimation creates a new Animation with frames and default 200ms speed
func NewAnimation(frames []string) Animation {
	return Animation{
		Frames: frames,
		speed:  200 * time.Millisecond,
	}
}

// Init returns initial tick command to start animation timer
func (a Animation) Init() Cmd {
	return Tick(a.speed)
}

// Update advances animation frame on TickMsg and returns next tick command
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

// View returns current animation frame as string for rendering
func (a Animation) View() string {
	if len(a.Frames) == 0 {
		return "Animation has no frames."
	}
	return a.Frames[a.frame]
}
