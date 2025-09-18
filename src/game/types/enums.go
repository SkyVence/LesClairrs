package types

type BodyParts int

const (
	Head BodyParts = iota
	Body
	ArmsL
	ArmsR
	Legs
)

type ItemType int

const (
	Upgrade ItemType = iota
	Utility
	Consumable
	Boost
	Weapon
	Armor
)

type CombatState int

const (
	OutOfCombat CombatState = iota
	Idle
	Attacking
	Defending
	Stunned
	Dead
	Victory
	EnemyTurn
	PlayerTurn
)
