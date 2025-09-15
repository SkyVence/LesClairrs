// Applique tous les bonus d'implants sur le joueur
func ApplyImplantBonuses(player *Player) {
    baseForce := 0
    baseSpeed := 0
    baseDefense := 0
    baseAccuracy := 0


    for _, implant := range player.Implants {
        for stat, bonus := range implant.StatBonus {
            switch stat {
            case "Force":
                baseForce += bonus
            case "Speed":
                baseSpeed += bonus
            case "Defense":
                baseDefense += bonus
            case "Accuracy":
                baseAccuracy += bonus
            }
        }
    }
    player.Force = baseForce
    player.Speed = baseSpeed
    player.Defense = baseDefense
    player.Accuracy = baseAccuracy
}
// Equipement d'un implant : un seul par slot
type PlayerImplants struct {
    Tete   *Implant
    BrasD  *Implant
    BrasG  *Implant
    Jambes *Implant
    Torse  *Implant
}

// Ajoute ou remplace l'implant du slot correspondant
func (pi *PlayerImplants) EquipImplant(implant *Implant) {
    switch implant.Slot {
    case "tete":
        pi.Tete = implant
    case "brasD":
        pi.BrasD = implant
    case "brasG":
        pi.BrasG = implant
    case "jambes":
        pi.Jambes = implant
    case "torse":
        pi.Torse = implant
    }
}
package main

import "fmt"

// Implant struct
type Implant struct {
    Name        string
    Slot        string
    StatBonus   map[string]int
    Skill       string
    Description string
    Cooldown    int // nombre de tours avant de pouvoir réutiliser le skill
    SkillEffect func(player *Player, enemy *Enemy)
}

lang,err = Load("fr")0
lang.Text()
var AllImplants = []Implant{
    // TETE
    {
        Name:      "Implant Cérébral I",
        Slot:      "tete",
        StatBonus: map[string]int{"Intelligence": 5},
        Skill:     "",
        Description: "Petit bonus d'intelligence.",
        Cooldown: 0,
        SkillEffect: nil,
    },
    {
        Name:      "Implant Cérébral II - Dissuasion",
        Slot:      "tete",
        StatBonus: map[string]int{"Intelligence": 10},
        Skill:     "dissuasion",
        Description: "Bloque le prochain coup de l'ennemi. (CD: 3 tours)",
        Cooldown: 3,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "bloque le prochain coup de l'ennemi !")
        },
    },
    {
        Name:      "Implant Cérébral III - Double Objet",
        Slot:      "tete",
        StatBonus: map[string]int{"Intelligence": 15},
        Skill:     "double_objet",
        Description: "Permet d'utiliser 2 objets en 1 tour. (CD: 4 tours)",
        Cooldown: 4,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "peut utiliser 2 objets en 1 tour !")
        },
    },

    // BRAS DROIT
    {
        Name:      "Implant de Force I",
        Slot:      "brasD",
        StatBonus: map[string]int{"Force": 8},
        Skill:     "",
        Description: "Petit bonus de force.",
        Cooldown: 0,
        SkillEffect: nil,
    },
    {
        Name:      "Implant de Force II - Smash",
        Slot:      "brasD",
        StatBonus: map[string]int{"Force": 18},
        Skill:     "smash",
        Description: "Smash : brise 30% de l'armure de l'ennemi. (CD: 3 tours)",
        Cooldown: 3,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Smash : brise 30% de l'armure de l'ennemi !")
        },
    },
    {
        Name:      "Implant de Force III - Tazer",
        Slot:      "brasD",
        StatBonus: map[string]int{"Force": 30},
        Skill:     "tazer",
        Description: "Tazer : inflige 25 dégâts et paralyse 1 tour. (CD: 4 tours)",
        Cooldown: 4,
        SkillEffect: func(player *Player, enemy *Enemy) {
            enemy.CurrentHP -= 25
            fmt.Println(player.Name, "utilise Tazer : inflige 25 dégâts et paralyse l'ennemi 1 tour !")
        },
    },

    // BRAS GAUCHE
    {
        Name:      "Implant de Précision I",
        Slot:      "brasG",
        StatBonus: map[string]int{"Accuracy": 8},
        Skill:     "",
        Description: "Petit bonus de précision.",
        Cooldown: 0,
        SkillEffect: nil,
    },
    {
        Name:      "Implant de Précision II - Piratage",
        Slot:      "brasG",
        StatBonus: map[string]int{"Accuracy": 18},
        Skill:     "piratage",
        Description: "Piratage : 50% de chance de skip le tour de l'ennemi. (CD: 3 tours)",
        Cooldown: 3,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Piratage : 50% de chance de skip le tour de l'ennemi !")
        },
    },
    {
        Name:      "Implant de Précision III - Headshot",
        Slot:      "brasG",
        StatBonus: map[string]int{"Accuracy": 30},
        Skill:     "headshot",
        Description: "Headshot : prochain tir x2 dégâts. (CD: 4 tours)",
        Cooldown: 4,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Headshot : prochain tir x2 dégâts !")
        },
    },

    // JAMBES
    {
        Name:      "Implant de Vitesse I",
        Slot:      "jambes",
        StatBonus: map[string]int{"Speed": 8},
        Skill:     "",
        Description: "Petit bonus de vitesse.",
        Cooldown: 0,
        SkillEffect: nil,
    },
    {
        Name:      "Implant de Vitesse II - Dash",
        Slot:      "jambes",
        StatBonus: map[string]int{"Speed": 18},
        Skill:     "dash",
        Description: "Dash : 50% de chances d'esquiver la prochaine attaque. (CD: 3 tours)",
        Cooldown: 3,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Dash : 50% de chances d'esquiver la prochaine attaque !")
        },
    },
    {
        Name:      "Implant de Vitesse III - Dash Parfait",
        Slot:      "jambes",
        StatBonus: map[string]int{"Speed": 30},
        Skill:     "dash_parfait",
        Description: "Dash Parfait : 100% d'esquive la prochaine attaque. (CD: 4 tours)",
        Cooldown: 4,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Dash Parfait : 100% d'esquive la prochaine attaque !")
        },
    },

    // TORSE
    {
        Name:      "Implant de Défense I",
        Slot:      "torse",
        StatBonus: map[string]int{"Defense": 8},
        Skill:     "",
        Description: "Petit bonus de défense.",
        Cooldown: 0,
        SkillEffect: nil,
    },
    {
        Name:      "Implant de Défense II - Torse en Fer",
        Slot:      "torse",
        StatBonus: map[string]int{"Defense": 18},
        Skill:     "torse_en_fer",
        Description: "Torse en Fer : réduit tous les dégâts entrants de 20% pendant 2 tours. (CD: 4 tours)",
        Cooldown: 4,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "utilise Torse en Fer : réduit tous les dégâts entrants de 20% pendant 2 tours !")
        },
    },
    {
        Name:      "Implant de Défense III - Bouclier Énergétique",
        Slot:      "torse",
        StatBonus: map[string]int{"Defense": 30},
        Skill:     "bouclier",
        Description: "Bouclier : absorbe tous les dégâts du prochain tour. (CD: 5 tours)",
        Cooldown: 5,
        SkillEffect: func(player *Player, enemy *Enemy) {
            fmt.Println(player.Name, "active Bouclier : absorbe tous les dégâts du prochain tour !")
        },
    },
}
