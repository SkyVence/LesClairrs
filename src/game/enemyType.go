package game

type Enemy struct {
	Name      string
	Force     int
	Speed     int
	Defense   int
	Accuracy  int
	MaxHP     int
	CurrentHP int
	Inventory []Item
	Implants  [5]Implant // "tete", "brasD", etc - fixed size array
	ExpReward int
}
