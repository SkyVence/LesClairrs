package main

import "fmt"


type Weapon struct {
    Name    string
    Type    string
    Attacks []Attack
}

type Attack struct {
    Name        string
    Description string
    Damage      int
    Duration    int // nombre de tours pour les effets continus
    Cooldown    int // nombre de tours avant de pouvoir réutiliser l'attaque
    Effect      func(player *Player, enemy *Enemy)
}

type Player struct {
    Name      string
    CurrentHP int
}

type Enemy struct {
    Name      string
    CurrentHP int
}


var AllWeapons = []Weapon{
    // KATANA
    {
        Name: "Katana",
        Type: "Melee",
        Attacks: []Attack{
            {"Frappe Spectrale", "Coup rapide, dégâts standards", 16, 0, 0, func(p *Player, e *Enemy) {
                e.CurrentHP -= 16
                fmt.Println(p.Name, "utilise Frappe Spectrale et inflige 16 dégâts à", e.Name)
            }},
            {"Slash Holo", "Dégât de saignement sur 2 tours", 12, 2, 2, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "inflige saignement à", e.Name, "pendant 2 tours")
            }},
            {"Imbattable", "Évite le prochain coup", 0, 3, 3, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "devient imbattable et esquive le prochain coup")
            }},
            {"ULTI – Lame Éthérée", "One shot mob, très gros dégâts boss", 140, 0, 3, func(p *Player, e *Enemy) {
                e.CurrentHP -= 140
                fmt.Println(p.Name, "utilise Lame Éthérée et inflige 140 dégâts à", e.Name)
            }},
        },
    },

    // FAUX NEURALINK
    {
        Name: "Faux Neuralink",
        Type: "Melee",
        Attacks: []Attack{
            {"Entaille", "Coup classique", 14, 0, 0, func(p *Player, e *Enemy) {
                e.CurrentHP -= 14
                fmt.Println(p.Name, "utilise Entaille et inflige 14 dégâts à", e.Name)
            }},
            {"Moisson Mentale", "Dégâts + vol de PV sur 2 tours", 10, 2, 2, func(p *Player, e *Enemy) {
                e.CurrentHP -= 10
                p.CurrentHP += 5 // vol de PV
                fmt.Println(p.Name, "utilise Moisson Mentale : 10 dégâts et récupère 5 PV")
            }},
            {"Rythme Instinctif", "Attaque spéciale", 0, 0, 2, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "utilise Rythme Instinctif")
            }},
            {"ULTI – Danse Macabre", "Dégâts massifs", 130, 0, 3, func(p *Player, e *Enemy) {
                e.CurrentHP -= 130
                fmt.Println(p.Name, "utilise Danse Macabre et inflige 130 dégâts à", e.Name)
            }},
        },
    },

    // ARC SYNAPTIQUE
    {
        Name: "Arc Synaptique",
        Type: "Ranged",
        Attacks: []Attack{
            {"Carreau EM", "Tir de base, dégâts électriques", 18, 0, 0, func(p *Player, e *Enemy) {
                e.CurrentHP -= 18
                fmt.Println(p.Name, "utilise Carreau EM et inflige 18 dégâts à", e.Name)
            }},
            {"Lien Conducteur", "Dégâts électriques pendant 3 tours", 7, 3, 2, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "utilise Lien Conducteur : 7 dégâts/round pendant 3 tours")
            }},
            {"Flèche Paralysante", "30% chance de skip le tour de l'ennemi", 16, 3, 3, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "utilise Flèche Paralysante : 30% de chance de skip le tour de l'ennemi")
            }},
            {"ULTI – Tempête Ionique", "Plusieurs grosses frappes électriques", 125, 0, 3, func(p *Player, e *Enemy) {
                e.CurrentHP -= 125
                fmt.Println(p.Name, "utilise Tempête Ionique et inflige 125 dégâts à", e.Name)
            }},
        },
    },

    // SN1P3R
    {
        Name: "SN1P3R",
        Type: "Ranged",
        Attacks: []Attack{
            {"Tir Balistique", "Dégâts élevés sur une seule cible", 22, 2, 2, func(p *Player, e *Enemy) {
                e.CurrentHP -= 22
                fmt.Println(p.Name, "utilise Tir Balistique et inflige 22 dégâts à", e.Name)
            }},
            {"Coup de Grâce", "Tir critique si ennemi <30% PV", 0, 3, 3, func(p *Player, e *Enemy) {
                if e.CurrentHP < e.CurrentHP/3 {
                    e.CurrentHP -= 60
                    fmt.Println(p.Name, "utilise Coup de Grâce : tir critique ! 60 dégâts")
                }
            }},
            {"Marqueur Optique", "Recharge une capacité gratuitement", 0, 2, 2, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "utilise Marqueur Optique et recharge une capacité")
            }},
            {"ULTI – H34DSH0T Parfait", "Ignorer l’armure et ×3 dégâts", 150, 0, 3, func(p *Player, e *Enemy) {
                e.CurrentHP -= 150
                fmt.Println(p.Name, "utilise H34DSH0T Parfait et inflige 150 dégâts à", e.Name)
            }},
        },
    },

    // NEON REAVER
    {
        Name: "Neon Reaver",
        Type: "Ranged",
        Attacks: []Attack{
            {"Tir Photonique", "Coup normal", 17, 0, 0, func(p *Player, e *Enemy) {
                e.CurrentHP -= 17
                fmt.Println(p.Name, "utilise Tir Photonique et inflige 17 dégâts à", e.Name)
            }},
            {"CHARGER!!", "Charge la barre de l'ultimate", 0, 0, 2, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "charge la barre de l'ULTI !")
            }},
            {"GRENADE!", "Explosion sur 3 tours", 10, 3, 3, func(p *Player, e *Enemy) {
                fmt.Println(p.Name, "lance GRENADE! : 10 dégâts/round pendant 3 tours")
            }},
            {"ULTI – Vortex Implosif", "Explosion massive", 130, 0, 3, func(p *Player, e *Enemy) {
                e.CurrentHP -= 130
                fmt.Println(p.Name, "utilise Vortex Implosif et inflige 130 dégâts à", e.Name)
            }},
        },
    },
}
