package types

type Position struct {
	X int
	Y int
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

// GetStage returns a pointer to the stage with the given number.
func (w *World) GetStage(stageNb int) *Stage {
	for i := range w.Stages {
		if w.Stages[i].StageNb == stageNb {
			return &w.Stages[i]
		}
	}
	return nil
}

// AdvanceToNextStage returns the next stage in the current world
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
// Note: This should probably use the loaders package to get the world
func (w *World) NextWorld() *World {
	if w == nil {
		return nil
	}
	// Return nil for now - this will be properly implemented with loaders
	return nil
}
