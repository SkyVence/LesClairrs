package game

import (
	"fmt"
)

// Structures simplifiées sans Ebiten
type Position struct {
	X, Y int
}

type Class struct {
	Name string
}

type Player struct {
	Name     string
	Class    Class
	Position Position
}

type Enemy struct {
	Name string
}

type Stage struct {
	WorldID        int
	StageNb        int
	Name           string
	Enemies        []Enemy
	ClearingReward int
}

type World struct {
	WorldID        int
	Name           string
	Stages         []Stage
	ClearingReward int
}

type Game struct {
	Player       *Player
	CurrentWorld *World
	CurrentStage *Stage
	dialogue     *DialogueBox
}

func NewPlayer(name string, class Class, pos Position) *Player {
	return &Player{
		Name:     name,
		Class:    class,
		Position: pos,
	}
}

func NewWorld(WorldID int) *World {
	return &World{
		WorldID: WorldID,
		Stages:  []Stage{}, // Initialiser avec des stages vides pour éviter les erreurs
	}
}

func NewGameInstance(selectedClass Class) *Game {
	world := NewWorld(1)

	// Ajouter au moins un stage pour éviter les erreurs
	world.Stages = append(world.Stages, Stage{
		WorldID:        1,
		StageNb:        1,
		Name:           "Premier Stage",
		Enemies:        []Enemy{},
		ClearingReward: 100,
	})

	player := NewPlayer("Sam", selectedClass, Position{X: 1, Y: 1})

	dialogue := NewDialogueBox()
	err := dialogue.LoadDialogues("fr.json") // ou le chemin vers votre fichier JSON
	if err != nil {
		fmt.Printf("Erreur chargement dialogues: %v\n", err)
	}

	return &Game{
		Player:       player,
		CurrentWorld: world,
		CurrentStage: &world.Stages[0],
		dialogue:     dialogue,
	}
}

// GetStage returns a pointer to the stage with the given number.
func (w *World) GetStage(stageNb int) *Stage {
	for i := range w.Stages {
		if w.Stages[i].StageNb == stageNb {
			return &w.Stages[i]
		}
	}
	return nil
}

func (s *Stage) AdvanceToNextStage(w *World) *Stage {
	if s == nil || w == nil {
		return nil
	}
	nextStageNb := s.StageNb + 1
	nextStage := w.GetStage(nextStageNb)
	if nextStage != nil {
		return nextStage
	}
	return nil
}

// NextWorld loads and returns the next world if available.
// It does not mutate the receiver; callers should assign the returned pointer.
func (w *World) NextWorld() *World {
	if w == nil {
		return nil
	}
	nextWorldID := w.WorldID + 1
	nextWorld := NewWorld(nextWorldID)
	if nextWorld != nil && len(nextWorld.Stages) > 0 {
		return nextWorld
	}
	return nil
}

// CurrentLocation returns a human-readable location string for the HUD.
func (g *Game) CurrentLocation() (string, int) {
	if g == nil || g.CurrentWorld == nil {
		return "", -1
	}
	if g.CurrentStage != nil && g.CurrentStage.Name != "" {
		// Prefer stage name when available
		return g.CurrentStage.Name, g.CurrentWorld.WorldID
	}
	return g.CurrentWorld.Name, g.CurrentWorld.WorldID
}

// PeekNext returns information about the next stage or world without mutating state.
// It returns a user-facing name for the upcoming location and the target world ID,
// plus a boolean indicating whether a next location exists.
func (g *Game) PeekNext() (string, int, bool) {
	if g == nil || g.CurrentWorld == nil || g.CurrentStage == nil {
		return "", -1, false
	}
	// Next stage in current world?
	if next := g.CurrentStage.AdvanceToNextStage(g.CurrentWorld); next != nil {
		name := next.Name
		if name == "" {
			name = g.CurrentWorld.Name
		}
		return name, g.CurrentWorld.WorldID, true
	}
	// Otherwise next world, first stage
	if nw := g.CurrentWorld.NextWorld(); nw != nil {
		if len(nw.Stages) > 0 {
			name := nw.Stages[0].Name
			if name == "" {
				name = nw.Name
			}
			return name, nw.WorldID, true
		}
	}
	return "", -1, false
}

// Advance moves the game to the next stage; if no more stages exist,
// it loads the next world and moves to its first stage. Returns true on success.
func (g *Game) Advance() bool {
	if g == nil || g.CurrentWorld == nil || g.CurrentStage == nil {
		return false
	}
	if next := g.CurrentStage.AdvanceToNextStage(g.CurrentWorld); next != nil {
		g.CurrentStage = next
		return true
	}
	// Try next world
	if nw := g.CurrentWorld.NextWorld(); nw != nil {
		g.CurrentWorld = nw
		if len(nw.Stages) > 0 {
			g.CurrentStage = &nw.Stages[0]
			return true
		}
	}
	return false
}

// Méthodes pour le dialogue terminal UNIQUEMENT
func (g *Game) StartDialogue() {
	if g.dialogue != nil {
		g.dialogue.StartDialogue()
	}
}

func (g *Game) ShowDialogueMenu() {
	fmt.Println("\n=== SYSTÈME DE DIALOGUE ===")
	fmt.Println("1. Démarrer les dialogues")
	fmt.Println("2. Retour au jeu")
	fmt.Print("Votre choix: ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		g.StartDialogue()
	case 2:
		return
	default:
		fmt.Println("Choix invalide.")
	}
}
