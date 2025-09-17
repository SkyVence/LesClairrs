package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type DialogueData struct {
	Dialogues []string `json:"dialogues"`
}

type DialogueStyles struct {
	Box       lipgloss.Style
	Text      lipgloss.Style
	Speaker   lipgloss.Style
	Indicator lipgloss.Style
}

func DefaultDialogueStyles() DialogueStyles {
	return DialogueStyles{
		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1a1a2e")).
			Padding(1, 2).
			Width(70).
			Height(8),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Width(66),
		Speaker: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8")).
			Background(lipgloss.Color("#654EA3")).
			Padding(0, 1).
			MarginBottom(1),
		Indicator: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8")),
	}
}

type DialogueBox struct {
	dialogues     []string
	currentIndex  int
	currentText   string
	displayedText string
	isVisible     bool
	styles        DialogueStyles
}

func NewDialogueBox() *DialogueBox {
	return &DialogueBox{
		dialogues: []string{},
		styles:    DefaultDialogueStyles(),
		isVisible: false,
	}
}

func (d *DialogueBox) LoadDialogues(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("erreur lecture fichier: %v", err)
	}

	var dialogueData DialogueData
	err = json.Unmarshal(data, &dialogueData)
	if err != nil {
		return fmt.Errorf("erreur parsing JSON: %v", err)
	}

	d.dialogues = dialogueData.Dialogues
	return nil
}

func (d *DialogueBox) StartDialogue() {
	if len(d.dialogues) == 0 {
		fmt.Println("Aucun dialogue à afficher.")
		return
	}

	d.isVisible = true
	d.currentIndex = 0
	d.showCurrentDialogue()
}

func (d *DialogueBox) showCurrentDialogue() {
	if d.currentIndex >= len(d.dialogues) {
		d.isVisible = false
		fmt.Println("\n" + d.styles.Indicator.Render("✓ Fin des dialogues"))
		return
	}

	d.currentText = d.dialogues[d.currentIndex]
	d.displayedText = d.currentText

	// Nettoyer l'écran et afficher le dialogue
	fmt.Print("\033[2J\033[H") // Clear screen
	d.render()

	// Attendre l'entrée utilisateur
	fmt.Print("\nAppuyez sur Entrée pour continuer...")
	fmt.Scanln()

	// Passer au suivant
	d.currentIndex++
	d.showCurrentDialogue()
}

func (d *DialogueBox) render() {
	var content []string

	// Extraire le speaker
	speaker, text := d.extractSpeaker(d.displayedText)

	if speaker != "" {
		content = append(content, d.styles.Speaker.Render(speaker))
		content = append(content, "")
	}

	// Découper le texte
	wrappedText := d.wrapText(text, 64)
	for _, line := range wrappedText {
		content = append(content, d.styles.Text.Render(line))
	}

	content = append(content, "")
	content = append(content, d.styles.Indicator.Render("▶ Appuyez sur Entrée pour continuer"))

	// Créer la boîte
	dialogueContent := lipgloss.JoinVertical(lipgloss.Left, content...)
	box := d.styles.Box.Render(dialogueContent)

	// Centrer sur l'écran
	width, height := 80, 24 // Taille standard terminal
	positioned := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)

	fmt.Print(positioned)
}

func (d *DialogueBox) extractSpeaker(text string) (string, string) {
	if colonIndex := strings.Index(text, ":"); colonIndex != -1 && colonIndex < 30 {
		speaker := strings.TrimSpace(text[:colonIndex])
		remainingText := strings.TrimSpace(text[colonIndex+1:])
		return speaker, remainingText
	}
	return "", text
}

func (d *DialogueBox) wrapText(text string, maxWidth int) []string {
	if text == "" {
		return []string{""}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func (d *DialogueBox) IsVisible() bool {
	return d.isVisible
}

func (d *DialogueBox) Close() {
	d.isVisible = false
}

// Fonction simple pour tester
func TestDialogue() {
	dialogue := NewDialogueBox()
	err := dialogue.LoadDialogues("fr.json")
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		return
	}

	dialogue.StartDialogue()
}