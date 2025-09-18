package types

type EnemyStats struct {
	Force     int
	Speed     int
	Defense   int
	Accuracy  int
	MaxHP     int
	CurrentHP int
}

type EnemySpawn struct {
	Name      string
	Force     int
	Speed     int
	Defense   int
	Accuracy  int
	MaxHP     int
	CurrentHP int
	Position  Position
	ExpReward int
	Sprite    string
}
