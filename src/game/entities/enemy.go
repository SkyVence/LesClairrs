package entities

import (
	"math"

	"projectred-rpg.com/game/types"
)

type Enemy struct {
	EnemyState types.CombatState
	Name       string
	Force      int
	Speed      int
	Defense    int
	Accuracy   int
	MaxHP      int
	CurrentHP  int
	ExpReward  int
	Sprite     string
	Position   types.Position
	IsAlive    bool
}

func NewEnemy(e Enemy) *Enemy {
	return &Enemy{
		Name:      e.Name,
		Force:     e.Force,
		Speed:     e.Speed,
		Defense:   e.Defense,
		Accuracy:  e.Accuracy,
		MaxHP:     e.MaxHP,
		CurrentHP: e.CurrentHP,
		ExpReward: e.ExpReward,
		Sprite:    e.Sprite,
		Position:  e.Position,
		IsAlive:   true,
	}
}

func (e *Enemy) GetStats() types.EnemyStats {
	return types.EnemyStats{
		Force:     e.Force,
		Speed:     e.Speed,
		Defense:   e.Defense,
		Accuracy:  e.Accuracy,
		MaxHP:     e.MaxHP,
		CurrentHP: e.CurrentHP,
	}
}

// TakeDamage applies damage to the enemy and returns true if the enemy is defeated
func (e *Enemy) TakeDamage(damage int) bool {
	e.CurrentHP -= damage
	if e.CurrentHP < 0 {
		e.CurrentHP = 0
		e.IsAlive = false
		return true
	}
	return false
}

func (e *Enemy) isDefeated() bool {
	return e.CurrentHP <= 0
}

// CalculateDamage computes the damage dealt to a target based on its defense
func (e *Enemy) CalculateDamage(targetDefense int) int {
	damage := e.Force - targetDefense
	if damage < 0 {
		damage = 0
	}
	return damage
}

func (e *Enemy) SetSprite(sprite string) {
	e.Sprite = sprite
}

// GetPosition returns the current position of the enemy
func (e *Enemy) GetPosition() types.Position {
	return e.Position
}

// SetPosition updates the enemy's position
func (e *Enemy) SetPosition(pos types.Position) {
	e.Position = pos
}

// DistanceTo calculates the distance between this enemy and a given position
func (e *Enemy) DistanceTo(pos types.Position) float64 {
	dx := float64(e.Position.X - pos.X)
	dy := float64(e.Position.Y - pos.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// IsWithinRange checks if the enemy is within a specified range of a position
func (e *Enemy) IsWithinRange(pos types.Position, maxDistance float64) bool {
	return e.DistanceTo(pos) <= maxDistance
}
