# 🎮 Game Developer API Documentation

A comprehensive guide for developing games using the ProjectRed RPG framework.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Game Package](#game-package)
3. [Player System](#player-system)
4. [World & Stage Management](#world--stage-management)
5. [Combat System](#combat-system)
6. [Inventory System](#inventory-system)
7. [Movement System](#movement-system)
8. [Entity Management](#entity-management)
9. [Asset Loading](#asset-loading)
10. [Configuration](#configuration)
11. [Complete Examples](#complete-examples)
12. [Best Practices](#best-practices)

---

## Overview

The ProjectRed RPG framework provides a complete set of game systems for building terminal-based RPGs. This documentation covers the high-level game development APIs that handle gameplay logic, character management, world progression, and game systems.

### Key Systems

- **Game Coordination**: Central game state management
- **Player Management**: Character stats, progression, and behavior
- **World System**: Stage-based progression and world loading
- **Combat System**: Battle mechanics and damage calculations
- **Inventory System**: Item management and equipment
- **Movement System**: Character movement and collision detection

---

## Game Package

### Core Game Coordinator

The `Game` struct serves as the central coordinator for all game systems:

```go
type Game struct {
    Player       *types.Player        // Current player instance
    CurrentWorld *types.World         // Active world
    CurrentStage *types.Stage         // Current stage within world
    
    // Game systems
    Combat    *systems.CombatSystem    // Combat calculations
    Inventory *systems.InventorySystem // Inventory management
    Movement  *systems.MovementSystem  // Movement handling
}
```

### Creating a New Game

```go
func NewGameInstance(selectedClass types.Class) *Game
```

Creates a fully initialized game instance with all systems ready.

**Parameters:**
- `selectedClass`: Character class determining base stats and abilities

**Example:**
```go
// Define a character class
class := types.Class{
    Name:        "Cyber-Samurai",
    Description: "A warrior enhanced with cybernetic implants",
    MaxHP:       100,
    Force:       15,
    Speed:       12,
    Defense:     8,
    Accuracy:    14,
}

// Create game instance
game := game.NewGameInstance(class)
```

### Game Navigation

#### Current Location Information

```go
func (g *Game) CurrentLocation() (string, int)
```

Returns human-readable location name and world ID for HUD display.

```go
location, worldID := game.CurrentLocation()
fmt.Printf("Current Location: %s (World %d)", location, worldID)
```

#### Peek Next Location

```go
func (g *Game) PeekNext() (string, int, bool)
```

Preview the next stage/world without advancing.

```go
nextLocation, nextWorldID, hasNext := game.PeekNext()
if hasNext {
    fmt.Printf("Next: %s (World %d)", nextLocation, nextWorldID)
}
```

#### Advance Progress

```go
func (g *Game) Advance() bool
```

Move to the next stage or world. Returns `true` on success.

```go
if game.Advance() {
    fmt.Println("Advanced to next stage!")
    location, _ := game.CurrentLocation()
    fmt.Printf("Now at: %s", location)
} else {
    fmt.Println("Cannot advance further")
}
```

---

## Player System

### Player Structure

```go
type Player struct {
    Name      string       // Player's name
    Class     Class        // Character class
    Stats     PlayerStats  // Current stats and progression
    Inventory []Item       // Player's items
    Position  Position     // Current position in world
    Implants  []Implant    // Cybernetic enhancements
}
```

### Character Classes

```go
type Class struct {
    Name        string  // Class name
    Description string  // Class description
    MaxHP       int     // Base health points
    Force       int     // Attack power
    Speed       int     // Movement and initiative
    Defense     int     // Damage reduction
    Accuracy    int     // Hit chance modifier
}
```

### Player Statistics

```go
type PlayerStats struct {
    Level        int     // Current level
    Exp          float32 // Current experience
    NextLevelExp int     // Experience needed for next level
    Force        int     // Current attack power
    Speed        int     // Current speed
    Defense      int     // Current defense
    Accuracy     int     // Current accuracy
    CurrentHP    int     // Current health
    MaxHP        int     // Maximum health
}
```

### Player Methods

#### Movement

```go
func (p *Player) Move(direction rune, worldWidth, worldHeight int) bool
```

Move player in the specified direction with boundary checking.

```go
// Handle player movement
func handleMovement(player *types.Player, key rune) {
    moved := player.Move(key, worldWidth, worldHeight)
    if moved {
        fmt.Println("Player moved successfully")
    } else {
        fmt.Println("Cannot move in that direction")
    }
}
```

#### Inventory Management

```go
func (p *Player) AddItemToInventory(item Item) bool
func (p *Player) RemoveItemFromInventory(itemName string) bool
func (p *Player) HasItem(itemName string) bool
```

**Examples:**
```go
// Add item to inventory
item := types.Item{
    Name:        "Health Potion",
    Type:        types.Consumable,
    Description: "Restores 50 HP",
    Value:       25,
}

if player.AddItemToInventory(item) {
    fmt.Println("Item added to inventory")
} else {
    fmt.Println("Inventory full!")
}

// Check for item
if player.HasItem("Health Potion") {
    fmt.Println("Player has health potion")
}

// Remove item
if player.RemoveItemFromInventory("Health Potion") {
    fmt.Println("Used health potion")
}
```

### Creating Players

```go
func NewPlayer(name string, class Class, position Position) *Player
```

**Example:**
```go
startPosition := types.Position{X: 1, Y: 1}
player := entities.NewPlayer("Sam", selectedClass, startPosition)
```

---

## World & Stage Management

### World Structure

```go
type World struct {
    WorldID     int      // Unique world identifier
    Name        string   // World name
    Description string   // World description
    Stages      []Stage  // Stages within this world
}
```

### Stage Structure

```go
type Stage struct {
    StageID     int      // Unique stage identifier
    Name        string   // Stage name
    Description string   // Stage description
    Width       int      // Stage width
    Height      int      // Stage height
    Entities    []Entity // NPCs, enemies, objects
}
```

### World Navigation Methods

```go
func (w *World) NextWorld() *World
func (s *Stage) AdvanceToNextStage(currentWorld *World) *Stage
```

### Creating Worlds

```go
func NewWorld(worldID int) *types.World
```

Loads world data from assets or creates empty world.

**Example:**
```go
// Load or create world
world := game.NewWorld(1)
fmt.Printf("Loaded world: %s", world.Name)

// Access stages
if len(world.Stages) > 0 {
    firstStage := world.Stages[0]
    fmt.Printf("First stage: %s", firstStage.Name)
}
```

---

## Combat System

### Combat System Structure

```go
type CombatSystem struct {
    // Combat configuration and state
}

func NewCombatSystem() *CombatSystem
```

### Combat Methods

#### Damage Calculation

```go
func (cs *CombatSystem) CalculateDamage(attackerForce, targetDefense int) int
```

Basic damage calculation with defense reduction.

#### Player Combat Actions

```go
func (cs *CombatSystem) PlayerAttacksEnemy(player *types.Player, enemy *types.Enemy) int
func (cs *CombatSystem) EnemyAttacksPlayer(enemy *types.Enemy, player *types.Player) int
```

#### Combat State Checks

```go
func (cs *CombatSystem) IsPlayerDefeated(player *types.Player) bool
func (cs *CombatSystem) IsEnemyDefeated(enemy *types.Enemy) bool
```

### Combat Example

```go
// Initialize combat
combat := systems.NewCombatSystem()

// Combat loop
for !combat.IsPlayerDefeated(player) && !combat.IsEnemyDefeated(enemy) {
    // Player turn
    damage := combat.PlayerAttacksEnemy(player, enemy)
    fmt.Printf("Player deals %d damage!", damage)
    
    if combat.IsEnemyDefeated(enemy) {
        fmt.Println("Enemy defeated!")
        // Award experience, loot, etc.
        break
    }
    
    // Enemy turn
    damage = combat.EnemyAttacksPlayer(enemy, player)
    fmt.Printf("Enemy deals %d damage!", damage)
    
    if combat.IsPlayerDefeated(player) {
        fmt.Println("Player defeated!")
        break
    }
}
```

---

## Inventory System

### Inventory System Structure

```go
type InventorySystem struct {
    MaxSlots int // Maximum inventory capacity
}

func NewInventorySystem() *InventorySystem
```

### Item Types

```go
type Item struct {
    Name        string   // Item name
    Type        ItemType // Item category
    Description string   // Item description
    Value       int      // Item value
    Rarity      Rarity   // Item rarity
    Effect      string   // Item effect description
}

type ItemType int
const (
    Weapon ItemType = iota
    Armor
    Consumable
    Accessory
    KeyItem
)
```

### Inventory Methods

```go
func (is *InventorySystem) AddItem(player *types.Player, item types.Item) bool
func (is *InventorySystem) RemoveItem(player *types.Player, itemName string) bool
func (is *InventorySystem) UseItem(player *types.Player, itemName string) bool
func (is *InventorySystem) GetInventoryCount(player *types.Player) int
func (is *InventorySystem) IsInventoryFull(player *types.Player) bool
```

### Inventory Example

```go
inventory := systems.NewInventorySystem()

// Create items
healthPotion := types.Item{
    Name:        "Health Potion",
    Type:        types.Consumable,
    Description: "Restores 50 HP",
    Value:       25,
    Effect:      "heal:50",
}

// Manage inventory
if !inventory.IsInventoryFull(player) {
    if inventory.AddItem(player, healthPotion) {
        fmt.Println("Added health potion to inventory")
    }
}

// Use item
if inventory.UseItem(player, "Health Potion") {
    fmt.Println("Used health potion - HP restored!")
}
```

---

## Movement System

### Movement System Structure

```go
type MovementSystem struct {
    // Movement configuration
}

func NewMovementSystem() *MovementSystem
```

### Movement Methods

```go
func (ms *MovementSystem) MovePlayer(player *types.Player, direction rune, world *types.World) bool
func (ms *MovementSystem) CheckCollision(position types.Position, world *types.World) bool
func (ms *MovementSystem) GetValidMoves(player *types.Player, world *types.World) []rune
```

### Position Structure

```go
type Position struct {
    X int // X coordinate
    Y int // Y coordinate
}
```

### Movement Example

```go
movement := systems.NewMovementSystem()

// Check valid moves
validMoves := movement.GetValidMoves(player, currentWorld)
fmt.Printf("Valid moves: %v", validMoves)

// Move player
direction := '↑' // Up arrow
if movement.MovePlayer(player, direction, currentWorld) {
    fmt.Println("Player moved successfully")
} else {
    fmt.Println("Movement blocked")
}
```

---

## Entity Management

### Player Entity Creation

```go
func NewPlayer(name string, class types.Class, position types.Position) *types.Player
```

### World Entity Creation

```go
func NewWorld(worldID int) *types.World
func GetWorld(worldID int) (types.World, bool)
```

### Entity Examples

```go
// Create player entity
startPos := types.Position{X: 5, Y: 5}
player := entities.NewPlayer("Hero", selectedClass, startPos)

// Create or load world entity
world := entities.NewWorld(1)
```

---

## Asset Loading

### World Loading

```go
func LoadWorlds() error
func GetWorld(worldID int) (types.World, bool)
```

### Asset Loading Examples

```go
// Load all world data
if err := loaders.LoadWorlds(); err != nil {
    log.Printf("Failed to load worlds: %v", err)
    // Handle gracefully with default content
}

// Get specific world
if world, exists := loaders.GetWorld(1); exists {
    fmt.Printf("Loaded world: %s", world.Name)
} else {
    fmt.Println("World not found")
}
```

---

## Configuration

### Game Constants

```go
// Game balance constants
const (
    MaxPlayerLevel = 50
    BaseExpPerLevel = 100
    MaxInventorySlots = 20
)

// Default character classes
var DefaultClasses = map[string]types.Class{
    "CYBER_SAMURAI": {
        Name:        "Cyber-Samurai",
        Description: "A warrior enhanced with cybernetic implants",
        MaxHP:       100,
        Force:       15,
        Speed:       12,
        Defense:     8,
        Accuracy:    14,
    },
    "TECH_MAGE": {
        Name:        "Tech-Mage",
        Description: "A spellcaster who manipulates technology",
        MaxHP:       80,
        Force:       10,
        Speed:       10,
        Defense:     5,
        Accuracy:    18,
    },
}
```

### Asset Paths

```go
const (
    WorldDataPath = "assets/levels/"
    AnimationPath = "assets/animations/"
    InterfacePath = "assets/interface/"
    DataPath      = "assets/data/"
)
```

---

## Complete Examples

### Simple Game Setup

```go
package main

import (
    "fmt"
    "log"
    "projectred-rpg.com/config"
    "projectred-rpg.com/game"
    "projectred-rpg.com/game/loaders"
)

func main() {
    // Load game assets
    if err := loaders.LoadWorlds(); err != nil {
        log.Printf("Warning: Failed to load worlds: %v", err)
    }
    
    // Select character class
    class := config.DefaultClasses["CYBER_SAMURAI"]
    
    // Create game instance
    gameInstance := game.NewGameInstance(class)
    
    // Display initial state
    location, worldID := gameInstance.CurrentLocation()
    fmt.Printf("Starting location: %s (World %d)\n", location, worldID)
    
    // Game loop example
    for {
        // Get user input (simplified)
        var input string
        fmt.Print("Enter command (move/advance/quit): ")
        fmt.Scanln(&input)
        
        switch input {
        case "advance":
            if gameInstance.Advance() {
                location, _ := gameInstance.CurrentLocation()
                fmt.Printf("Advanced to: %s\n", location)
            } else {
                fmt.Println("Cannot advance further")
            }
        case "quit":
            return
        default:
            fmt.Println("Unknown command")
        }
    }
}
```

### Combat System Example

```go
func combatExample(player *types.Player, enemy *types.Enemy) {
    combat := systems.NewCombatSystem()
    
    fmt.Printf("Combat started: %s vs %s\n", player.Name, enemy.Name)
    
    round := 1
    for !combat.IsPlayerDefeated(player) && !combat.IsEnemyDefeated(enemy) {
        fmt.Printf("\n--- Round %d ---\n", round)
        
        // Player turn
        damage := combat.PlayerAttacksEnemy(player, enemy)
        fmt.Printf("%s attacks for %d damage!\n", player.Name, damage)
        fmt.Printf("%s HP: %d/%d\n", enemy.Name, enemy.CurrentHP, enemy.MaxHP)
        
        if combat.IsEnemyDefeated(enemy) {
            fmt.Printf("%s defeated!\n", enemy.Name)
            // Award experience, loot, etc.
            break
        }
        
        // Enemy turn
        damage = combat.EnemyAttacksPlayer(enemy, player)
        fmt.Printf("%s attacks for %d damage!\n", enemy.Name, damage)
        fmt.Printf("%s HP: %d/%d\n", player.Name, player.Stats.CurrentHP, player.Stats.MaxHP)
        
        if combat.IsPlayerDefeated(player) {
            fmt.Printf("%s defeated!\n", player.Name)
            break
        }
        
        round++
    }
}
```

### Inventory Management Example

```go
func inventoryExample(player *types.Player) {
    inventory := systems.NewInventorySystem()
    
    // Create some items
    items := []types.Item{
        {
            Name:        "Health Potion",
            Type:        types.Consumable,
            Description: "Restores 50 HP",
            Value:       25,
        },
        {
            Name:        "Cyber Sword",
            Type:        types.Weapon,
            Description: "A high-tech blade",
            Value:       150,
        },
    }
    
    // Add items to inventory
    for _, item := range items {
        if inventory.AddItem(player, item) {
            fmt.Printf("Added %s to inventory\n", item.Name)
        } else {
            fmt.Printf("Could not add %s - inventory full\n", item.Name)
        }
    }
    
    // Check inventory status
    count := inventory.GetInventoryCount(player)
    fmt.Printf("Inventory: %d items\n", count)
    
    // Use an item
    if inventory.UseItem(player, "Health Potion") {
        fmt.Println("Used health potion")
    }
}
```

---

## Best Practices

### 1. Game State Management

- Keep game state centralized in the `Game` struct
- Use the systems architecture for modular functionality
- Handle state transitions explicitly

```go
// Good: Clear state management
func (g *Game) AdvanceToNextArea() bool {
    if g.canAdvance() {
        g.performAdvancement()
        g.updateGameState()
        return true
    }
    return false
}
```

### 2. Error Handling

- Handle missing assets gracefully
- Provide fallback content when possible
- Log errors without crashing

```go
// Good: Graceful error handling
if err := loaders.LoadWorlds(); err != nil {
    log.Printf("Warning: Could not load worlds: %v", err)
    // Continue with default/empty world
    world := &types.World{WorldID: 1, Name: "Default World"}
}
```

### 3. Player Progression

- Validate player actions before applying changes
- Maintain consistency in stat calculations
- Save progression state appropriately

```go
// Good: Validated progression
func (p *Player) GainExperience(amount float32) bool {
    if amount <= 0 {
        return false
    }
    
    p.Stats.Exp += amount
    if p.Stats.Exp >= float32(p.Stats.NextLevelExp) {
        return p.levelUp()
    }
    return true
}
```

### 4. System Integration

- Use dependency injection for system coupling
- Keep systems loosely coupled
- Provide clear interfaces between systems

```go
// Good: Loose coupling
type GameSystems struct {
    Combat    CombatInterface
    Inventory InventoryInterface
    Movement  MovementInterface
}
```

### 5. Asset Management

- Organize assets in clear directory structures
- Implement caching for frequently accessed data
- Handle missing assets gracefully

### 6. Configuration Management

- Use constants for game balance values
- Make configuration easily modifiable
- Provide sensible defaults

---

This documentation provides comprehensive coverage of the ProjectRed RPG game development framework, enabling developers to create rich terminal-based RPG experiences using the provided systems and APIs.