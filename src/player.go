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
