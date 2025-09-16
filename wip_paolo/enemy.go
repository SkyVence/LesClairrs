package main

import (
    "fmt"
    "math/rand"
    "time"
)


type EnemyAttack struct {
    Name        string
    Damage      int
    Cooldown    int // nombre de tours avant de pouvoir réutiliser
    Description string
    Effect      func(player *Player, enemy *Enemy)
}

type Enemy struct {
    Name       string
    MaxHP      int
    CurrentHP  int
    Attacks    []EnemyAttack
    Type       string
}


var AllEnemies = []Enemy{
    // Mob normal : 1 attaque simple
    {
        Name: "Mob Standard",
        MaxHP: 50,
        CurrentHP: 50,
        Attacks: []EnemyAttack{
            {Name: "Coup de griffe", Damage: 10, Cooldown: 0, Description: "Attaque de base", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 10
            }},
        },
        Type: "normal",
    },
    // Mini boss : 2 attaques
    {
        Name: "Mini Boss",
        MaxHP: 120,
        CurrentHP: 120,
        Attacks: []EnemyAttack{
            {Name: "Coup puissant", Damage: 18, Cooldown: 0, Description: "Attaque forte", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 18
            }},
            {Name: "Hurlement", Damage: 12, Cooldown: 2, Description: "Effraie le joueur, utilisable tous les 2 tours", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 12
            }},
        },
        Type: "mini boss",
    },
    // Boss : 3 attaques
    {
        Name: "Boss",
        MaxHP: 300,
        CurrentHP: 300,
        Attacks: []EnemyAttack{
            {Name: "Écrasement", Damage: 28, Cooldown: 0, Description: "Attaque principale", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 28
            }},
            {Name: "Onde de choc", Damage: 20, Cooldown: 2, Description: "Attaque de zone, utilisable tous les 2 tours", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 20
            }},
            {Name: "Rugissement", Damage: 15, Cooldown: 3, Description: "Affaiblit le joueur, utilisable tous les 3 tours", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 15
            }},
        },
        Type: "boss",
    },
    // Boss de fin de jeu : 3 attaques puissantes
    {
        Name: "Boss Final",
        MaxHP: 500,
        CurrentHP: 500,
        Attacks: []EnemyAttack{
            {Name: "Frappe du Néant", Damage: 40, Cooldown: 0, Description: "Attaque ultime", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 40
            }},
            {Name: "Explosion d'énergie", Damage: 30, Cooldown: 2, Description: "Explosion, utilisable tous les 2 tours", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 30
            }},
            {Name: "Malédiction", Damage: 25, Cooldown: 3, Description: "Malus sur le joueur, utilisable tous les 3 tours", Effect: func(p *Player, e *Enemy) {
                p.CurrentHP -= 25
            }},
        },
        Type: "boss de fin de jeu",
    },
}



func GetDamageFromPlayer(player *Player, enemy *Enemy, attack Attack) {
    rand.Seed(time.Now().UnixNano())
    hitChance := rand.Intn(100)
    if hitChance < player.Accuracy {
        baseDamage := attack.Damage + player.Force - enemy.Defense
        if baseDamage < 0 {
            baseDamage = 0
        }
        enemy.CurrentHP -= baseDamage
        fmt.Printf("%s utilise %s et inflige %d dégâts à %s !\n", player.Name, attack.Name, baseDamage, enemy.Name)
        if attack.Effect != nil {
            attack.Effect(player, enemy)
        }
    } else {
        fmt.Printf("%s utilise %s mais rate son attaque !\n", player.Name, attack.Name)
    }
}


func DoDamageToPlayer(player *Player, enemy *Enemy) {
    // Choisir une attaque (ici la première pour l'exemple)
    if len(enemy.Attacks) == 0 {
        return
    }
    attack := enemy.Attacks[0]
    // Application des dégâts
    damage := attack.Damage
    player.CurrentHP -= damage
    if player.CurrentHP < 0 {
        player.CurrentHP = 0
    }
    fmt.Printf("%s utilise %s et inflige %d dégâts à %s !\n", enemy.Name, attack.Name, damage, player.Name)
    // Application de l'effet spécial si présent
    if attack.Effect != nil {
        attack.Effect(player, enemy)
    }
}
