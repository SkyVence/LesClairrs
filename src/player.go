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
