package systems

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"projectred-rpg.com/engine"
	"projectred-rpg.com/ui"
)

// LevelIntroSystem handles level introduction dialogues
type LevelIntroSystem struct {
	dialogBox          *ui.DialogBox
	localization       map[string]interface{}
	language           string
	isActive           bool
	onComplete         func()
	currentDialogIndex int
	currentDialogs     []DialogLine
}

type DialogLine struct {
	Speaker string
	Message string
	Order   int
}

// NewLevelIntroSystem creates a new level intro system
func NewLevelIntroSystem(language string) *LevelIntroSystem {
	return &LevelIntroSystem{
		dialogBox:      ui.NewDialogBox(100), // Ajoutez la largeur par défaut
		language:       language,
		isActive:       false,
		currentDialogs: make([]DialogLine, 0),
	}
}

// LoadLocalization loads localization data from JSON files
func (lis *LevelIntroSystem) LoadLocalization() error {
	filename := fmt.Sprintf("assets/interface/%s.json", lis.language)

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read localization file %s: %v", filename, err)
	}

	if err := json.Unmarshal(data, &lis.localization); err != nil {
		return fmt.Errorf("failed to parse localization file %s: %v", filename, err)
	}

	return nil
}

// ShowIntro displays the introduction for a specific level
func (lis *LevelIntroSystem) ShowIntro(levelFilename string, screenWidth, screenHeight int, onComplete func()) bool {
	// Extract world and stage from filename (e.g., "world-1_stage-1.map" -> world1, stage1)
	levelName := strings.TrimSuffix(levelFilename, ".map")
	parts := strings.Split(levelName, "_")

	fmt.Printf("DEBUG: levelFilename = %s\n", levelFilename)
	fmt.Printf("DEBUG: levelName = %s\n", levelName)
	fmt.Printf("DEBUG: parts = %v\n", parts)

	if len(parts) != 2 {
		fmt.Printf("DEBUG: parts length is not 2, got %d\n", len(parts))
		return false
	}

	worldPart := strings.Replace(parts[0], "world-", "world", 1)
	stagePart := strings.Replace(parts[1], "stage-", "", 1)

	fmt.Printf("DEBUG: worldPart = %s\n", worldPart)
	fmt.Printf("DEBUG: stagePart = %s\n", stagePart)

	// Navigate through the JSON structure
	game, ok := lis.localization["game"].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: 'game' not found in localization\n")
		return false
	}

	levels, ok := game["levels"].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: 'levels' not found in game\n")
		return false
	}

	world, ok := levels[worldPart].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: world '%s' not found in levels\n", worldPart)
		return false
	}

	stages, ok := world["stages"].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: 'stages' not found in world\n")
		return false
	}

	stage, ok := stages[stagePart].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: stage '%s' not found in stages\n", stagePart)
		return false
	}

	dialogue, ok := stage["dialogue"].(map[string]interface{})
	if !ok {
		fmt.Printf("DEBUG: 'dialogue' not found in stage\n")
		return false
	}

	fmt.Printf("DEBUG: Found dialogue with %d entries\n", len(dialogue))

	// Convert dialogue to DialogLine array
	lis.currentDialogs = lis.convertDialogueToLines(dialogue)
	if len(lis.currentDialogs) == 0 {
		fmt.Printf("DEBUG: No dialogs converted\n")
		return false
	}

	fmt.Printf("DEBUG: Converted %d dialogues, starting intro\n", len(lis.currentDialogs))

	lis.onComplete = onComplete
	lis.isActive = true
	lis.currentDialogIndex = 0

	// Show first dialog
	firstDialog := lis.currentDialogs[0]
	fmt.Printf("DEBUG: First dialog: %s says: %s\n", firstDialog.Speaker, firstDialog.Message)
	lis.dialogBox.ShowCentered(firstDialog.Message, firstDialog.Speaker, screenWidth, screenHeight)

	return true
}

// convertDialogueToLines converts the dialogue map to ordered dialog lines
func (lis *LevelIntroSystem) convertDialogueToLines(dialogue map[string]interface{}) []DialogLine {
	lines := make([]DialogLine, 0, len(dialogue))

	for key, value := range dialogue {
		if message, ok := value.(string); ok {
			speaker, order := lis.parseDialogueKey(key)
			lines = append(lines, DialogLine{
				Speaker: speaker,
				Message: message,
				Order:   order,
			})
		}
	}

	// Sort by order
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Order < lines[j].Order
	})

	return lines
}

// parseDialogueKey extracts speaker name and order from dialogue key
func (lis *LevelIntroSystem) parseDialogueKey(key string) (string, int) {
	// Find the last digit(s) in the key
	i := len(key) - 1
	for i >= 0 && key[i] >= '0' && key[i] <= '9' {
		i--
	}

	if i == len(key)-1 {
		// No number found, return order 0
		return key, 0
	}

	speaker := key[:i+1]
	orderStr := key[i+1:]

	order, err := strconv.Atoi(orderStr)
	if err != nil {
		return key, 0
	}

	return speaker, order
}

// Update updates the intro system
func (lis *LevelIntroSystem) Update(msg engine.Msg) (*LevelIntroSystem, engine.Cmd) {
	if !lis.isActive {
		return lis, nil
	}

	switch msg := msg.(type) {
	case engine.KeyMsg:
		switch msg.Rune {
		case '\r', ' ': // Enter key or space
			if lis.dialogBox.IsTextComplete() {
				// Move to next dialog
				lis.currentDialogIndex++
				if lis.currentDialogIndex >= len(lis.currentDialogs) {
					// All dialogs shown, complete intro
					lis.Hide()
					if lis.onComplete != nil {
						lis.onComplete()
					}
				} else {
					// Show next dialog
					nextDialog := lis.currentDialogs[lis.currentDialogIndex]
					lis.dialogBox.ShowCentered(nextDialog.Message, nextDialog.Speaker, 80, 24)
				}
			} else {
				// Complete current text immediately
				lis.dialogBox.AdvanceText()
			}
		}
	}

	// Update dialog box
	var cmd engine.Cmd
	lis.dialogBox, cmd = lis.dialogBox.Update(msg)

	return lis, cmd
}

// Render renders the intro dialog
func (lis *LevelIntroSystem) Render() string {
	if !lis.isActive {
		return ""
	}
	return lis.dialogBox.Render()
}

// IsActive returns whether an intro is currently being shown
func (lis *LevelIntroSystem) IsActive() bool {
	return lis.isActive
}

// Hide hides the intro dialog
func (lis *LevelIntroSystem) Hide() {
	lis.isActive = false
	lis.dialogBox.Hide()
}