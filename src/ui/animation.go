package ui

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

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
		trimmedFrame := strings.TrimSpace(frame)
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

// ViewAligned returns the current animation frame with consistent positioning.
// The maxWidth parameter ensures that all frames are aligned properly,
// which helps prevent visual "jumping" when frames have different widths.
func (a Animation) ViewAligned(maxWidth int) string {
	if len(a.Frames) == 0 {
		return "Animation has no frames."
	}

	// Create a style with a fixed width and center alignment to ensure consistent positioning
	style := lipgloss.NewStyle().Width(maxWidth).Align(lipgloss.Center)
	return style.Render(a.Frames[a.frame])
}
