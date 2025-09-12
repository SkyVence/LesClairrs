# Tutoriel Framework UI - Création de Composants et Animations

## Architecture du Framework

Ce framework UI suit le pattern **The Elm Architecture (TEA)** avec trois concepts principaux :

### 1. Model (Modèle)
Le modèle représente l'état de votre application. Il contient toutes les données nécessaires pour rendre l'interface utilisateur.

### 2. Update (Mise à jour)
La fonction `Update` gère les messages entrants et met à jour le modèle en conséquence. Elle retourne un nouveau modèle et optionnellement une commande.

### 3. View (Vue)
La fonction `View` prend le modèle actuel et retourne une représentation string de l'interface utilisateur.

## Cycle de Vie

```
┌─────────┐    ┌──────────┐    ┌─────────┐
│  Init   │───▶│  Update  │───▶│  View   │
└─────────┘    └────┬─────┘    └─────────┘
                    │                │
                    ▼                ▼
               ┌─────────┐      ┌─────────┐
               │   Cmd   │      │ Render  │
               └─────────┘      └─────────┘
```

## Création d'un Composant Simple

Voici comment créer un composant compteur basique :

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/votre-projet/ui"
)

// Modèle du compteur
type CounterModel struct {
    count int
}

// Messages que le compteur peut recevoir
type IncrementMsg struct{}
type DecrementMsg struct{}

// Implémentation de l'interface Model
func (m CounterModel) Init() ui.Msg {
    return nil // Pas de message initial
}

func (m CounterModel) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg := msg.(type) {
    case IncrementMsg:
        m.count++
        return m, nil
    
    case DecrementMsg:
        m.count--
        return m, nil
    
    case ui.KeyMsg:
        switch msg.Rune {
        case '+':
            return m, func() ui.Msg { return IncrementMsg{} }
        case '-':
            return m, func() ui.Msg { return DecrementMsg{} }
        case 'q':
            return m, func() ui.Msg { return ui.Quit() }
        }
    }
    
    return m, nil
}

func (m CounterModel) View() string {
    return fmt.Sprintf(`
╭─────────────────────╮
│     COMPTEUR        │
├─────────────────────┤
│                     │
│    Valeur: %3d      │
│                     │
│  [+] Incrémenter    │
│  [-] Décrémenter    │
│  [q] Quitter        │
│                     │
╰─────────────────────╯
`, m.count)
}
```

## Animations avec le Système Tick

Le framework fournit un système de tick pour créer des animations :

### Exemple : Horloge Animée

```go
package main

import (
    "fmt"
    "time"
    "github.com/votre-projet/ui"
)

type ClockModel struct {
    currentTime time.Time
}

func (m ClockModel) Init() ui.Msg {
    return ui.TickNow()() // Démarre immédiatement
}

func (m ClockModel) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg := msg.(type) {
    case ui.TickMsg:
        m.currentTime = msg.Time
        // Programme le prochain tick dans 1 seconde
        return m, ui.Tick(time.Second)
    
    case ui.KeyMsg:
        if msg.Rune == 'q' {
            return m, func() ui.Msg { return ui.Quit() }
        }
    }
    
    return m, nil
}

func (m ClockModel) View() string {
    timeStr := m.currentTime.Format("15:04:05")
    return fmt.Sprintf(`
╭─────────────────────╮
│      HORLOGE        │
├─────────────────────┤
│                     │
│    %s         │
│                     │
│   [q] Quitter       │
╰─────────────────────╯
`, timeStr)
}
```

### Exemple : Barre de Progression Animée

```go
package main

import (
    "fmt"
    "strings"
    "time"
    "github.com/votre-projet/ui"
)

type ProgressModel struct {
    progress int // 0-100
    direction int // 1 ou -1
}

func (m ProgressModel) Init() ui.Msg {
    return ui.TickNow()()
}

func (m ProgressModel) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg := msg.(type) {
    case ui.TickMsg:
        // Mettre à jour la progression
        m.progress += m.direction * 2
        
        // Inverser la direction aux limites
        if m.progress >= 100 {
            m.progress = 100
            m.direction = -1
        } else if m.progress <= 0 {
            m.progress = 0
            m.direction = 1
        }
        
        return m, ui.Tick(50 * time.Millisecond)
    
    case ui.KeyMsg:
        if msg.Rune == 'q' {
            return m, func() ui.Msg { return ui.Quit() }
        }
    }
    
    return m, nil
}

func (m ProgressModel) View() string {
    barWidth := 20
    filled := int(float64(barWidth) * float64(m.progress) / 100.0)
    
    bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
    
    return fmt.Sprintf(`
╭─────────────────────────╮
│    BARRE PROGRESSION    │
├─────────────────────────┤
│                         │
│ [%s] %3d%%  │
│                         │
│     [q] Quitter         │
╰─────────────────────────╯
`, bar, m.progress)
}
```

## Utilisation avec le Programme Principal

```go
func main() {
    // Créer le modèle initial
    model := CounterModel{count: 0}
    // ou ClockModel{currentTime: time.Now()}
    // ou ProgressModel{progress: 0, direction: 1}
    
    // Créer le programme
    program := ui.NewProgram(model, ui.WithAltScreen())
    
    // Lancer l'application
    if err := program.Run(); err != nil {
        fmt.Printf("Erreur: %v\n", err)
    }
}
```

## Messages Personnalisés

Vous pouvez définir vos propres types de messages :

```go
// Messages pour un jeu
type MoveLeftMsg struct{}
type MoveRightMsg struct{}
type JumpMsg struct{}
type GameOverMsg struct {
    Score int
}

// Messages pour des événements réseau
type DataReceivedMsg struct {
    Data []byte
}

type ConnectionLostMsg struct {
    Error error
}
```

## Bonnes Pratiques

1. **Gardez le modèle simple** : Ne stockez que l'état nécessaire
2. **Messages spécifiques** : Créez des messages clairs et précis
3. **Commandes asynchrones** : Utilisez des goroutines pour les opérations longues
4. **Gestion d'erreurs** : Incluez l'état d'erreur dans votre modèle
5. **Performance** : Évitez les recalculs coûteux dans View()

## Exemple Complet : Mini-Jeu

```go
type GameModel struct {
    playerX   int
    enemyX    int
    score     int
    gameOver  bool
}

func (m GameModel) Init() ui.Msg {
    return ui.TickNow()()
}

func (m GameModel) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    if m.gameOver {
        if key, ok := msg.(ui.KeyMsg); ok && key.Rune == 'r' {
            // Redémarrer le jeu
            return GameModel{playerX: 10, enemyX: 0, score: 0}, ui.TickNow()
        }
        return m, nil
    }
    
    switch msg := msg.(type) {
    case ui.TickMsg:
        // Logique du jeu
        m.enemyX++
        if m.enemyX > 20 {
            m.enemyX = 0
            m.score++
        }
        
        // Collision ?
        if abs(m.playerX - m.enemyX) < 2 {
            m.gameOver = true
            return m, nil
        }
        
        return m, ui.Tick(200 * time.Millisecond)
    
    case ui.KeyMsg:
        switch msg.Rune {
        case 'a':
            if m.playerX > 0 {
                m.playerX--
            }
        case 'd':
            if m.playerX < 20 {
                m.playerX++
            }
        case 'q':
            return m, func() ui.Msg { return ui.Quit() }
        }
    }
    
    return m, nil
}

func (m GameModel) View() string {
    if m.gameOver {
        return fmt.Sprintf(`
GAME OVER!
Score: %d

[r] Redémarrer
[q] Quitter
`, m.score)
    }
    
    // Créer la zone de jeu
    field := make([]rune, 21)
    for i := range field {
        field[i] = '.'
    }
    field[m.playerX] = 'P'
    field[m.enemyX] = 'E'
    
    return fmt.Sprintf(`
Score: %d

%s

[a] Gauche [d] Droite [q] Quitter
`, m.score, string(field))
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
```

Ce framework vous permet de créer des applications interactives complexes en combinant modèles, messages et animations de manière élégante et prévisible.
