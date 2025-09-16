package game

type Game struct {
	Player *Player
	// Level  *Level // TODO: Implement level loading and management
}

func NewGameInstance() *Game {
	// Hardcoded for now
	player := NewPlayer("Hero", Class{
		Name:     "Cyber-Samurai",
		Force:    10,
		Speed:    12,
		Defense:  8,
		Accuracy: 15,
		MaxHP:    100,
	}, Position{X: 10, Y: 5})

	return &Game{
		Player: player,
	}
}
