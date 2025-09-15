package main

import (
	"fmt"
)

// Prix des items (à ajuster selon l'équilibrage)
var ItemPrices = map[string]int{
	"Bandages douteux": 10,
	"HealPack": 25,
	"TON1C": 30,
	"M0NE¥": 30,
	"TOX1C": 20,
	"flashbang": 20,
	"ADR3N4LIN3": 20,
}

// Affiche la liste des items disponibles à l'achat
func ShowShop() {
	fmt.Println("--- Marchand ---")
	for i, item := range AllItems {
		price := ItemPrices[item.Name]
		fmt.Printf("%d. %s (%s) - %d$\n   %s\n", i+1, item.Name, item.Rarity, price, item.Description)
	}
}

// Permet au joueur d'acheter un item (retourne true si achat réussi)
func BuyItem(player *Player, itemIndex int) bool {
	if itemIndex < 0 || itemIndex >= len(AllItems) {
		fmt.Println("Choix invalide, fais un effort !")
		return false
	}
	item := AllItems[itemIndex]
	price := ItemPrices[item.Name]
	if player.Money < price {
		fmt.Println("Pas assez d'argent ! Retourne me voir quand t'as les moyens.")
		return false
	}
	// Ajoute / incrémente l'item dans l'inventaire
	found := false
	for i, invItem := range player.Inventory {
		if invItem.Item.Name == item.Name {
			player.Inventory[i].Quantity++
			found = true
			break
		}
	}
	if !found {
		player.Inventory = append(player.Inventory, InventoryItem{Item: item, Quantity: 1})
	}
	player.Money -= price
	fmt.Printf("%s acheté pour %d$ !\n", item.Name, price)
	return true
}
