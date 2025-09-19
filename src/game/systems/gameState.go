package systems

type StateEnum int

const (
	StateMainMenu StateEnum = iota
	StateSettings
	StateClassSelection
	StateExploration
	StateCombat
	StateMerchant
	StateDialogue
	StateInventory
	StateDeathScreen
	StateVictoryScreen
	StatePauseMenu
	StateStageTransition
)

type GameState struct {
	CurrentState  StateEnum
	PreviousState StateEnum
}

func NewGameState(initial StateEnum) *GameState {
	return &GameState{
		CurrentState:  initial,
		PreviousState: initial,
	}
}

func (gs *GameState) ChangeState(newState StateEnum) {
	gs.PreviousState = gs.CurrentState
	gs.CurrentState = newState
}

func (gs *GameState) RevertState() {
	gs.CurrentState, gs.PreviousState = gs.PreviousState, gs.CurrentState
}

func (gs *GameState) IsInState(state StateEnum) bool {
	return gs.CurrentState == state
}
