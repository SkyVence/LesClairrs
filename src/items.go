// Objet dans l'inventaire du joueur (avec quantité)
type InventoryItem struct {
	Item     Item
	Quantity int
}
package main

import (
    "fmt"
    "math/rand"
    "time"
)


type Item struct {
	Name        string
	Type        string
	Rarity      string
	Description string
	Cooldown    int // nombre de tours avant de pouvoir réutiliser l'objet (si applicable)
	Effect      func(player *Player, enemy *Enemy)
}
var AllItems = []Item{
	{
		Name: "Bandages douteux",
		Type: "soin",
		Rarity: "commun",
		Description: "Rend 25% des PV max.",
		Cooldown: 0,
		Effect: func(p *Player, e *Enemy) {
			heal := p.MaxHP / 4
			p.CurrentHP += heal
			if p.CurrentHP > p.MaxHP {
				p.CurrentHP = p.MaxHP
			}
			fmt.Println(p.Name, "récupère", heal, "PV !")
		},
	},
	{
		Name:   "HealPack",
		Type:   "soin",
		Rarity: "rare",
		Description: "Rend 50% des PV max.",
		Cooldown: 0,
		Effect: func(p *Player, e *Enemy) {
			heal := p.MaxHP / 2
			p.CurrentHP += heal
			if p.CurrentHP > p.MaxHP {
				p.CurrentHP = p.MaxHP
			}
			fmt.Println(p.Name, "récupère", heal, "PV !")
		},
	},
	{
		Name:   "TON1C",
		Type:   "boost",
		Rarity: "rare",
		Description: "Double l'expérience gagnée ce niveau. (CD: 1 combat)",
		Cooldown: 1,
		Effect: func(p *Player, e *Enemy) {
			fmt.Println(p.Name, "double l'expérience pour ce niveau !")
		},
	},
	{
		Name:   "M0NE¥",
		Type:   "boost",
		Rarity: "rare",
		Description: "Double l'argent gagné ce niveau. (CD: 1 combat)",
		Cooldown: 1,
		Effect: func(p *Player, e *Enemy) {
			fmt.Println(p.Name, "double l'argent pour ce niveau !")
		},
	},
	{
		Name:   "TOX1C",
		Type:   "utility",
		Rarity: "commun",
		Description: "Empoisonne l'ennemi ciblé pour 3 tours (5 dégâts/tour). (CD: 2 tours)",
		Cooldown: 2,
		Effect: func(p *Player, e *Enemy) {
			if e != nil {
				if e.Status == nil {
					e.Status = make(map[string]int)
				}
				e.Status["poison"] = 3
				e.Status["poison_dmg"] = 5
				fmt.Println(p.Name, "empoisonne", e.Name, "pour 3 tours !")
			}
		},
	},
	{
		Name:   "flashbang",
		Type:   "utility",
		Rarity: "commun",
		Description: "Étourdît l'ennemi ciblé (skip son prochain tour). (CD: 3 tours)",
		Cooldown: 3,
		Effect: func(p *Player, e *Enemy) {
			if e != nil {
				if e.Status == nil {
					e.Status = make(map[string]int)
				}
				e.Status["stun"] = 1
				fmt.Println(p.Name, "étourdit", e.Name, "pour 1 tour !")
			}
		},
	},
	{
		Name:   "ADR3N4LIN3",
		Type:   "utility",
		Rarity: "commun",
		Description: "Permet d'attaquer 2 fois ce tour. (CD: 2 tours)",
		Cooldown: 2,
		Effect: func(p *Player, e *Enemy) {
			if p.Status == nil {
				p.Status = make(map[string]int)
			}
			p.Status["double_attack"] = 1
			fmt.Println(p.Name, "attaque 2 fois ce tour !")
		},
	},
	{
		Name:   "Grande Bourse",
		Type:   "upgrade",
		Rarity: "rare",
		Description: "Augmente la capacité d'inventaire de 10 slots.",
		Cooldown: 0,
		Effect: func(p *Player, e *Enemy) {
			p.MaxInv += 10
			fmt.Println(p.Name, "utilise Grande Bourse : capacité d'inventaire +10 !")
		},
	},
	{
		Name:   "Grosse Bourse",
		Type:   "upgrade",
		Rarity: "épique",
		Description: "Augmente la capacité d'inventaire de 10 slots.",
		Cooldown: 0,
		Effect: func(p *Player, e *Enemy) {
			p.MaxInv += 10
			fmt.Println(p.Name, "utilise Grosse Bourse : capacité d'inventaire +10 !")
		},
	},
}
