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