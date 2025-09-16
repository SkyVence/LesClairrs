package game

type bodyParts int
type rarity int
type ItemType int

const (
	Upgrade = iota
	Utility
	Consumable
	Boost
)

const (
	Head = iota
	Body
	ArmsL
	ArmsR
	Legs
)

const (
	Tier1 = iota
	Tier2
	Tier3
)

type Class struct {
	Name        string
	Description string
	MaxHP       int
	Force       int
	Speed       int
	Defense     int
	Accuracy    int
}

type BonusStats struct {
	Force    int
	Speed    int
	Defense  int
	Accuracy int
}

type Implant struct {
	Type        bodyParts
	Rarity      rarity
	Bonus       BonusStats
	Name        string
	Description string
	Cooldown    int
	Ability     ImplantAbility
}

type PlayerStats struct {
	Level        int
	Exp          float32
	NextLevelExp int
	Force        int
	Speed        int
	Defense      int
	Accuracy     int
	MaxHP        int
	CurrentHP    int
}

type Item struct {
	Name        string
	Type        ItemType
	Rarity      rarity
	Description string
	Cooldown    int // nombre de tours avant de pouvoir réutiliser l'objet (si applicable)
	Effect      Effect
}

type Position struct {
	X int
	Y int
}

type Player struct {
	Name  string
	Class Class
	Stats PlayerStats
	Pos   Position

	// Placeholder until i implement moving animation
	sprite string

	Inventory []Item
	Implants  [5]Implant // "tete", "brasD", etc - fixed size array
	MaxInv    int
}

type Effect interface {
	Use(user *Player, target *Enemy)
	Name() string
	CoolDown() int
}

type ImplantAbility interface {
	OnEquip(owner *Player)
	OnUnequip(owner *Player)
	Name() string
}
