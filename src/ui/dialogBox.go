package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
)

// DialogBox represents a dialog box UI component
type DialogBox struct {
	content     string
	speaker     string
	visible     bool
	x, y        int
	width       int
	height      int
	textIndex   int
	isComplete  bool
	showCursor  bool
	styles      DialogBoxStyles
}

// DialogBoxStyles contains styling for the dialog box
type DialogBoxStyles struct {
	Border      lipgloss.Style
	Content     lipgloss.Style
	Speaker     lipgloss.Style
	Background  lipgloss.Style
}

// DefaultDialogBoxStyles returns the default dialog box styling
func DefaultDialogBoxStyles() DialogBoxStyles {
	return DialogBoxStyles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2),
		Content: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")),
		Speaker: lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true),
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("0")),
	}
}

// NewDialogBox creates a new dialog box
func NewDialogBox(width int, styles ...DialogBoxStyles) *DialogBox {
	var style DialogBoxStyles
	if len(styles) > 0 {
		style = styles[0]
	} else {
		style = DefaultDialogBoxStyles()
	}

	return &DialogBox{
		width:  width,
		height: 8,
		styles: style,
	}
}

// Show displays the dialog box with the specified content and speaker
func (d *DialogBox) Show(content, speaker string, x, y int) {
	d.content = content
	d.speaker = speaker
	d.x = x
	d.y = y
	d.visible = true
	d.textIndex = 0
	d.isComplete = false
	d.showCursor = true
}

// ShowCentered displays the dialog box centered on screen
func (d *DialogBox) ShowCentered(content, speaker string, screenWidth, screenHeight int) {
	d.content = content
	d.speaker = speaker
	d.width = screenWidth - 4
	d.height = 8
	d.x = 2
	d.y = screenHeight - d.height - 2
	d.visible = true
	d.textIndex = 0
	d.isComplete = false
	d.showCursor = true
}

// Hide hides the dialog box
func (d *DialogBox) Hide() {
	d.visible = false
}

// IsVisible returns whether the dialog box is currently visible
func (d *DialogBox) IsVisible() bool {
	return d.visible
}

// IsTextComplete returns whether the text animation is complete
func (d *DialogBox) IsTextComplete() bool {
	return d.isComplete
}

// AdvanceText advances the typewriter effect or marks as complete
func (d *DialogBox) AdvanceText() {
	if d.textIndex < len(d.content) {
		d.textIndex = len(d.content)
		d.isComplete = true
	}
}

// Update handles messages and updates the dialog box state
func (d *DialogBox) Update(msg engine.Msg) (*DialogBox, engine.Cmd) {
	if !d.visible {
		return d, nil
	}

	switch msg := msg.(type) {
	case engine.KeyMsg:
		switch msg.Rune {
		case '\r', ' ': // Enter key or space
			if !d.isComplete {
				d.AdvanceText()
			}
		}
	case engine.TickMsg:
		// Typewriter effect
		if d.textIndex < len(d.content) && !d.isComplete {
			d.textIndex++
			if d.textIndex >= len(d.content) {
				d.isComplete = true
			}
		}
		// Cursor blinking
		d.showCursor = !d.showCursor
	}

	return d, nil
}

// wrapText wraps text to fit within the specified width
func (d *DialogBox) wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

// Render renders the dialog box to a string
func (d *DialogBox) Render() string {
	if !d.visible {
		return ""
	}

	// Get displayed text (with typewriter effect)
	displayText := d.content
	if d.textIndex < len(d.content) {
		displayText = d.content[:d.textIndex]
	}

	// Add cursor if text is complete and visible
	if d.isComplete && d.showCursor {
		displayText += " ▋"
	}

	// Wrap text
	wrappedText := d.wrapText(displayText, d.width-6)

	// Build dialog content
	var content strings.Builder
	if d.speaker != "" {
		content.WriteString(d.styles.Speaker.Render(d.speaker + ":"))
		content.WriteString("\n")
	}
	content.WriteString(d.styles.Content.Render(wrappedText))

	// Apply border and return
	return d.styles.Border.
		Width(d.width).
		Height(d.height).
		Render(content.String())
}