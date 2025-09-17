package types

type Item struct {
	Type        ItemType
	Rarity      Rarity
	Name        string
	Description string
}

type Weapon struct {
	KeyName string
	Type    int
	Attacks []Attack
}

type Attack struct {
	KeyName  string
	KeyDesc  string
	Damage   int
	Duration int
	CoolDown int
}
