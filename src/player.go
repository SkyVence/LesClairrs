type Class struct {
    Name        string
    Description string
    MaxHP       int
    Force       int
    Speed       int
    Defense     int
    Accuracy    int
}

var AllClasses = []Class{
    {
        Name: "D0C",
        Description: "Un robot intelligent, précis et polyvalent.",
        MaxHP: 90,
        Force: 10,
        Speed: 12,
        Defense: 10,
        Accuracy: 22,
    },
    {
        Name: "APP",
        Description: "Un robot furtif, rapide et précis.",
        MaxHP: 80,
        Force: 14,
        Speed: 22,
        Defense: 8,
        Accuracy: 18,
    },
    {
        Name: "PER",
        Description: "Un robot robuste, fort et endurant.",
        MaxHP: 120,
        Force: 22,
        Speed: 10,
        Defense: 18,
        Accuracy: 10,
    },
}

// Initialise un joueur selon la classe choisie
func NewPlayer(name string, className string) *Player {
    var chosen Class
    for _, c := range AllClasses {
        if c.Name == className {
            chosen = c
            break
        }
    }
    return &Player{
        Name: name,
        Class: chosen.Name,
        MaxHP: chosen.MaxHP,
        CurrentHP: chosen.MaxHP,
        Force: chosen.Force,
        Speed: chosen.Speed,
        Defense: chosen.Defense,
        Accuracy: chosen.Accuracy,
        Inventory: []InventoryItem{},
        Implants: make(map[string]Implant),
        MaxInv: 10,
    }
}
// Ajoute un objet à l'inventaire, en respectant la limite MaxInv
func (p *Player) AddItem(item Item, qty int) bool {
    total := 0
    for _, invItem := range p.Inventory {
        total += invItem.Quantity
    }
    if total+qty > p.MaxInv {
        fmt.Println("Inventaire plein !")
        return false
    }
    for i, invItem := range p.Inventory {
        if invItem.Item.Name == item.Name {
            p.Inventory[i].Quantity += qty
            return true
        }
    }
    p.Inventory = append(p.Inventory, InventoryItem{Item: item, Quantity: qty})
    return true
}
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Player struct {
    Name      string
    Class     string
    Level     int
    MaxHP     int
    CurrentHP int
    Force     int
    Speed     int
    Defense   int
    Accuracy  int
    Inventory []InventoryItem
    Implants  map[string]Implant // "tete", "brasD", etc
    MaxInv    int
}

// Utilise un objet de l'inventaire (usage unique)
func (p *Player) UseItem(itemName string, enemy *Enemy) bool {
    for i, invItem := range p.Inventory {
        if invItem.Item.Name == itemName && invItem.Quantity > 0 {
            invItem.Item.Effect(p, enemy)
            p.Inventory[i].Quantity--
            if p.Inventory[i].Quantity == 0 {
                // Retire l'objet de l'inventaire
                p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
            }
            return true
        }
    }
    return false
}
