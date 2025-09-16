package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"projectred-rpg.com/engine"
)

// HUD represents a heads-up display component that sticks to the bottom
type HUD struct {
	width      int
	height     int
	termWidth  int
	termHeight int

	// Game state to display
	playerHealth    int
	playerMaxHealth int
	playerLevel     int
	playerExp       int
	expToNextLevel  int
	location        string
}

// HUDStyles contains styling for the HUD
type HUDStyles struct {
	Container lipgloss.Style
	HealthBar lipgloss.Style
	ExpBar    lipgloss.Style
	Text      lipgloss.Style
	Border    lipgloss.Style
}

// DefaultHUDStyles returns the default HUD styling
func DefaultHUDStyles() HUDStyles {
	return HUDStyles{
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1),
		HealthBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
		ExpBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")),
		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")),
	}
}

// NewHud creates a new HUD component with default values
func NewHud() *HUD {
	return &HUD{
		width:           80, // Default width
		height:          4,  // Fixed height for HUD (3 lines + borders)
		playerHealth:    100,
		playerMaxHealth: 100,
		playerLevel:     1,
		playerExp:       0,
		expToNextLevel:  100,
		location:        "Starting Area",
	}
}

func (h *HUD) Init() engine.Cmd {
	return engine.TickNow()
}

func (h *HUD) Update(msg engine.Msg) (HUD, engine.Cmd) {
	switch msg := msg.(type) {
	case engine.SizeMsg:
		h.termWidth = msg.Width
		h.termHeight = msg.Height
		// HUD takes full width but fixed height
		h.width = msg.Width
	}
	return *h, nil
}

// SetPlayerStats updates the player statistics displayed in the HUD
func (h *HUD) SetPlayerStats(health, maxHealth, level, exp, expToNext int, location string) {
	h.playerHealth = health
	h.playerMaxHealth = maxHealth
	h.playerLevel = level
	h.playerExp = exp
	h.expToNextLevel = expToNext
	h.location = location
}

// View renders the HUD as a bottom-positioned component
func (h *HUD) View() string {
	if h.width == 0 || h.termWidth == 0 {
		return ""
	}

	styles := DefaultHUDStyles()

	// Calculate bar widths (accounting for padding and borders)
	availableWidth := h.width - 4 // Account for borders and padding
	if availableWidth < 30 {
		availableWidth = 30 // Minimum available width
	}
	barWidth := availableWidth / 3 // Split into 3 sections

	if barWidth < 10 {
		barWidth = 10 // Minimum bar width
	}
	if barWidth > 20 {
		barWidth = 20 // Maximum bar width for readability
	}

	// Create health bar
	healthPercent := float64(h.playerHealth) / float64(h.playerMaxHealth)
	healthFilled := int(float64(barWidth) * healthPercent)
	if healthFilled < 0 {
		healthFilled = 0
	}
	if healthFilled > barWidth {
		healthFilled = barWidth
	}
	healthBar := strings.Repeat("█", healthFilled) + strings.Repeat("▒", barWidth-healthFilled)

	// Create experience bar
	expPercent := float64(h.playerExp) / float64(h.expToNextLevel)
	expFilled := int(float64(barWidth) * expPercent)
	if expFilled < 0 {
		expFilled = 0
	}
	if expFilled > barWidth {
		expFilled = barWidth
	}
	expBar := strings.Repeat("█", expFilled) + strings.Repeat("▒", barWidth-expFilled)

	// Format the HUD content
	healthText := fmt.Sprintf("HP: %d/%d", h.playerHealth, h.playerMaxHealth)
	expText := fmt.Sprintf("EXP: %d/%d", h.playerExp, h.expToNextLevel)
	levelText := fmt.Sprintf("Level %d", h.playerLevel)
	locationText := h.location

	// Create the three sections
	leftSection := lipgloss.JoinVertical(lipgloss.Left,
		styles.Text.Render(healthText),
		styles.HealthBar.Render(healthBar),
	)

	centerSection := lipgloss.JoinVertical(lipgloss.Center,
		styles.Text.Render(levelText),
		styles.Text.Render(locationText),
	)

	rightSection := lipgloss.JoinVertical(lipgloss.Right,
		styles.Text.Render(expText),
		styles.ExpBar.Render(expBar),
	)

	// Join sections horizontally with proper spacing
	sectionWidth := (availableWidth - 4) / 3 // Account for spacing
	if sectionWidth < 10 {
		sectionWidth = 10 // Minimum section width
	}

	leftFormatted := lipgloss.NewStyle().Width(sectionWidth).Align(lipgloss.Left).Render(leftSection)
	centerFormatted := lipgloss.NewStyle().Width(sectionWidth).Align(lipgloss.Center).Render(centerSection)
	rightFormatted := lipgloss.NewStyle().Width(sectionWidth).Align(lipgloss.Right).Render(rightSection)

	hudContent := lipgloss.JoinHorizontal(lipgloss.Top,
		leftFormatted,
		" ", // Spacer
		centerFormatted,
		" ", // Spacer
		rightFormatted,
	)

	// Apply container styling with border
	styledHUD := styles.Container.Width(h.width - 2).Render(hudContent)

	return styledHUD
}

// RenderWithContent combines main content with the bottom-positioned HUD
// This is a helper function to properly position the HUD at the bottom
func (h *HUD) RenderWithContent(mainContent string) string {
	if h.termWidth == 0 || h.termHeight == 0 {
		return mainContent // Fallback if size not set
	}

	hudView := h.View()

	// Calculate available space for main content
	hudLines := strings.Count(hudView, "\n") + 1
	mainContentHeight := h.termHeight - hudLines

	// Ensure main content doesn't overflow
	if mainContentHeight < 1 {
		mainContentHeight = 1
	}

	// Position main content in the available space above the HUD
	positionedMainContent := lipgloss.Place(
		h.termWidth, mainContentHeight,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)

	// Position HUD at the bottom
	positionedHUD := lipgloss.Place(
		h.termWidth, hudLines,
		lipgloss.Center, lipgloss.Bottom,
		hudView,
	)

	// Combine them vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		positionedMainContent,
		positionedHUD,
	)
}
