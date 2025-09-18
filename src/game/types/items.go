package types

type Item struct {
	Type        ItemType
	Name        string
	Description string
}

// Renomme Weapon en WeaponData pour éviter le conflit
type WeaponData struct {
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
