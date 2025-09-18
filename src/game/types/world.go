package types

type Position struct {
	X int
	Y int
}

type Stage struct {
	WorldID        int
	StageNb        int
	Name           string
	Enemies        []EnemySpawn
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
