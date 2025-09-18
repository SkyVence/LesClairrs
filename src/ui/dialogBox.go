package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/types"
)

// DialogBox represents a dialog box UI component that can be positioned relative to entities
type DialogBox struct {
	width       int
	height      int
	maxWidth    int
	content     string
	speaker     string
	position    types.Position
	visible     bool
	styles      DialogBoxStyles
	textIndex   int  // For typewriter effect
	showCursor  bool // For blinking cursor
	isComplete  bool // Whether text animation is complete
}

// DialogBoxStyles contains styling for the dialog box
type DialogBoxStyles struct {
	Container    lipgloss.Style
	Text         lipgloss.Style
	Speaker      lipgloss.Style
	Border       lipgloss.Style
	Cursor       lipgloss.Style
	Background   lipgloss.Style
}

// DefaultDialogBoxStyles returns the default dialog box styling
func DefaultDialogBoxStyles() DialogBoxStyles {
	return DialogBoxStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#4A90E2")).
			Padding(1, 2).
			MarginBottom(1).
			Background(lipgloss.Color("#1A1A2E")),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Speaker: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD700")).
			MarginBottom(1),
		Border: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A90E2")),
		Cursor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true),
		Background: lipgloss.NewStyle().
			Background(lipgloss.Color("#1A1A2E")).
			Foreground(lipgloss.Color("#FFFFFF")),
	}
}

// NewDialogBox creates a new dialog box with the specified parameters
func NewDialogBox(maxWidth int, styles ...DialogBoxStyles) *DialogBox {
	dialogStyles := DefaultDialogBoxStyles()
	if len(styles) > 0 {
		dialogStyles = styles[0]
	}

	return &DialogBox{
		maxWidth:   maxWidth,
		visible:    false,
		styles:     dialogStyles,
		textIndex:  0,
		showCursor: true,
		isComplete: false,
	}
}

// Show displays the dialog box with the specified content and speaker
func (d *DialogBox) Show(content, speaker string, npcPos types.Position) {
	d.content = content
	d.speaker = speaker
	d.position = d.calculatePosition(npcPos)
	d.visible = true
	d.textIndex = 0
	d.isComplete = false
}

// Hide hides the dialog box
func (d *DialogBox) Hide() {
	d.visible = false
	d.content = ""
	d.speaker = ""
	d.textIndex = 0
	d.isComplete = false
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
		d.textIndex = len(d.content) // Skip to end
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

// calculatePosition determines where to position the dialog box relative to the NPC
func (d *DialogBox) calculatePosition(npcPos types.Position) types.Position {
	// Position the dialog box above and slightly to the right of the NPC
	return types.Position{
		X: npcPos.X + 2,
		Y: npcPos.Y - 5, // Above the NPC
	}
}

// wrapText wraps text to fit within the specified width
func (d *DialogBox) wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// If adding this word would exceed the width, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Render renders the dialog box to a string
func (d *DialogBox) Render() string {
	if !d.visible {
		return ""
	}

	// Calculate actual width for text wrapping (accounting for padding and borders)
	textWidth := d.maxWidth - 6 // Account for borders and padding

	// Get the text to display (with typewriter effect)
	displayText := d.content
	if d.textIndex < len(d.content) {
		displayText = d.content[:d.textIndex]
	}

	// Wrap the text
	wrappedLines := d.wrapText(displayText, textWidth)

	// Add cursor if text is not complete
	if !d.isComplete && d.showCursor && len(wrappedLines) > 0 {
		lastLineIndex := len(wrappedLines) - 1
		wrappedLines[lastLineIndex] += d.styles.Cursor.Render("▌")
	}

	// Build the content
	var content strings.Builder

	// Add speaker name if present
	if d.speaker != "" {
		content.WriteString(d.styles.Speaker.Render(d.speaker))
		content.WriteString("\n")
	}

	// Add wrapped text lines
	for i, line := range wrappedLines {
		content.WriteString(d.styles.Text.Render(line))
		if i < len(wrappedLines)-1 {
			content.WriteString("\n")
		}
	}

	// Add interaction hint if text is complete
	if d.isComplete {
		content.WriteString("\n")
		content.WriteString(d.styles.Cursor.Render("Press ENTER to continue..."))
	}

	// Apply container styling
	return d.styles.Container.Render(content.String())
}

// GetPosition returns the dialog box position
func (d *DialogBox) GetPosition() types.Position {
	return d.position
}

// SetPosition sets the dialog box position
func (d *DialogBox) SetPosition(pos types.Position) {
	d.position = pos
}

// GetDimensions returns the dialog box dimensions
func (d *DialogBox) GetDimensions() (int, int) {
	return d.width, d.height
}