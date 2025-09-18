package systems

import (
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/types"
	"projectred-rpg.com/ui"
)

// DialogEntry represents a single dialog entry with localization key and speaker
type DialogEntry struct {
	SpeakerKey string   // Localization key for speaker name
	TextKey    string   // Localization key for dialog text
	Args       []any    // Arguments for text formatting
}

// DialogSequence represents a sequence of dialog entries
type DialogSequence struct {
	Entries []DialogEntry
	Current int
}

// DialogSystem manages dialog interactions and UI
type DialogSystem struct {
	dialogBox     *ui.DialogBox
	isActive      bool
	currentDialog *DialogSequence
	locManager    *engine.LocalizationManager
	onComplete    func() // Callback when dialog sequence completes
}

// NewDialogSystem creates a new dialog system
func NewDialogSystem(maxWidth int) *DialogSystem {
	return &DialogSystem{
		dialogBox:  ui.NewDialogBox(maxWidth),
		isActive:   false,
		locManager: engine.GetLocalizationManager(),
	}
}

// StartDialog begins a new dialog sequence with the specified NPC
func (ds *DialogSystem) StartDialog(sequence *DialogSequence, npcPos types.Position, onComplete func()) {
	if sequence == nil || len(sequence.Entries) == 0 {
		return
	}

	ds.currentDialog = sequence
	ds.currentDialog.Current = 0
	ds.isActive = true
	ds.onComplete = onComplete

	// Display the first dialog entry
	ds.showCurrentEntry(npcPos)
}

// EndDialog ends the current dialog sequence
func (ds *DialogSystem) EndDialog() {
	ds.isActive = false
	ds.currentDialog = nil
	ds.dialogBox.Hide()
	
	if ds.onComplete != nil {
		ds.onComplete()
		ds.onComplete = nil
	}
}

// IsActive returns whether a dialog is currently active
func (ds *DialogSystem) IsActive() bool {
	return ds.isActive
}

// Update processes messages and updates the dialog system
func (ds *DialogSystem) Update(msg engine.Msg) {
	if !ds.isActive {
		return
	}

	// Update the dialog box
	ds.dialogBox, _ = ds.dialogBox.Update(msg)

	// Handle dialog progression
	if keyMsg, ok := msg.(engine.KeyMsg); ok {
		switch keyMsg.Rune {
		case '\r', ' ': // Enter or Space
			if ds.dialogBox.IsTextComplete() {
				ds.nextEntry()
			}
		case 27: // Escape key
			ds.EndDialog()
		}
	}
}

// Render returns the rendered dialog box
func (ds *DialogSystem) Render() string {
	if !ds.isActive {
		return ""
	}
	return ds.dialogBox.Render()
}

// GetDialogBox returns the dialog box for external positioning/rendering
func (ds *DialogSystem) GetDialogBox() *ui.DialogBox {
	return ds.dialogBox
}

// showCurrentEntry displays the current dialog entry
func (ds *DialogSystem) showCurrentEntry(npcPos types.Position) {
	if ds.currentDialog == nil || ds.currentDialog.Current >= len(ds.currentDialog.Entries) {
		ds.EndDialog()
		return
	}

	entry := ds.currentDialog.Entries[ds.currentDialog.Current]
	
	// Get localized text
	speakerText := ""
	if entry.SpeakerKey != "" {
		speakerText = ds.locManager.Text(entry.SpeakerKey)
	}
	
	dialogText := ds.locManager.Text(entry.TextKey, entry.Args...)
	
	// Show the dialog box
	ds.dialogBox.Show(dialogText, speakerText, npcPos.X, npcPos.Y)
}

// nextEntry advances to the next dialog entry or ends the dialog
func (ds *DialogSystem) nextEntry() {
    if ds.currentDialog == nil {
        return
    }

    ds.currentDialog.Current++
    
    if ds.currentDialog.Current >= len(ds.currentDialog.Entries) {
        // End of dialog sequence
        ds.EndDialog()
    } else {
        // Show next entry - use a default position since GetPosition doesn't exist
        defaultPos := types.Position{X: 2, Y: 5} // Position par défaut
        ds.showCurrentEntry(defaultPos)
    }
}

// CreateSimpleDialog creates a simple dialog sequence with a single entry
func CreateSimpleDialog(speakerKey, textKey string, args ...any) *DialogSequence {
	return &DialogSequence{
		Entries: []DialogEntry{
			{
				SpeakerKey: speakerKey,
				TextKey:    textKey,
				Args:       args,
			},
		},
		Current: 0,
	}
}

// CreateMultiDialog creates a dialog sequence with multiple entries
func CreateMultiDialog(entries ...DialogEntry) *DialogSequence {
	return &DialogSequence{
		Entries: entries,
		Current: 0,
	}
}

// AddEntry adds a new entry to an existing dialog sequence
func (ds *DialogSequence) AddEntry(speakerKey, textKey string, args ...any) {
	entry := DialogEntry{
		SpeakerKey: speakerKey,
		TextKey:    textKey,
		Args:       args,
	}
	ds.Entries = append(ds.Entries, entry)
}

// HasNext returns whether there are more entries in the sequence
func (ds *DialogSequence) HasNext() bool {
	return ds.Current < len(ds.Entries)-1
}

// Reset resets the dialog sequence to the beginning
func (ds *DialogSequence) Reset() {
	ds.Current = 0
}