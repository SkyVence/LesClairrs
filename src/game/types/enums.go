package types

type BodyParts int

const (
	Head BodyParts = iota
	Body
	ArmsL
	ArmsR
	Legs
)

type Rarity int

const (
	Tier1 Rarity = iota
	Tier2
	Tier3
)

type ItemType int

const (
	Upgrade ItemType = iota
	Utility
	Consumable
	Boost
)
