package entities

import "projectred-rpg.com/game/types"

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
