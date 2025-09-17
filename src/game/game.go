package game

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
}

func NewWorld(WorldID int) *World {
	// Ensure worlds are loaded
	if err := LoadWorlds(); err != nil {
		return &World{WorldID: WorldID}
	}

	if world, exists := GetWorld(WorldID); exists {
		return &world
	}
	return &World{WorldID: WorldID}
}

func NewGameInstance(selectedClass Class) *Game {

	world := NewWorld(1)

	// Hardcoded for now --> Probably will be changed if implementing save/load system
	player := NewPlayer("Sam", selectedClass, Position{X: 1, Y: 1})	
	return &Game{
		Player:       player,
		CurrentWorld: world,
		CurrentStage: &world.Stages[0],
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

	// Use cached world data instead of loading from file
	if nextWorld, exists := GetWorld(nextWorldID); exists {
		return &nextWorld
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
