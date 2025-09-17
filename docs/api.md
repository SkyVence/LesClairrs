# 📚 API Reference

This document provides comprehensive API documentation for ProjectRed RPG.

## 📋 Table of Contents

- [Game Package](#game-package)
- [Types Package](#types-package)
- [Systems Package](#systems-package)
- [Entities Package](#entities-package)
- [Loaders Package](#loaders-package)
- [Config Package](#config-package)

## 🎮 Game Package

### Core Game Coordinator

#### `type Game`
Main game state coordinator that integrates all systems and manages game flow.

```go
type Game struct {
    Player       *types.Player        // Current player instance
    CurrentWorld *types.World         // Active world
    CurrentStage *types.Stage         // Current stage within world
    
    // Game systems
    Combat    *systems.CombatSystem     // Combat calculations
    Inventory *systems.InventorySystem  // Inventory management
    Movement  *systems.MovementSystem   // Movement handling
}
```

#### `func NewGameInstance(selectedClass types.Class) *Game`
Creates a new game instance with the specified character class.

**Parameters:**
- `selectedClass`: The character class for the player

**Returns:** Configured Game instance with initialized systems

**Example:**
```go
class := types.Class{Name: "Cyber-Samurai", MaxHP: 100, Force: 15}
game := NewGameInstance(class)
```

#### `func (g *Game) CurrentLocation() (string, int)`
Returns the current location name and world ID for display.

**Returns:**
- `string`: Human-readable location name
- `int`: World ID

**Example:**
```go
location, worldID := game.CurrentLocation()
fmt.Printf("Current location: %s (World %d)", location, worldID)
```

#### `func (g *Game) PeekNext() (string, int, bool)`
Previews the next stage/world without advancing the game state.

**Returns:**
- `string`: Name of next location
- `int`: Target world ID
- `bool`: Whether a next location exists

#### `func (g *Game) Advance() bool`
Advances to the next stage or world.

**Returns:** `true` if advancement was successful

## 🏗️ Types Package

### Player Types

#### `type Player`
Represents a player character with stats, inventory, and behavior.

```go
type Player struct {
    Name      string           // Player name
    Class     Class           // Character class
    Stats     PlayerStats     // Current statistics
    Pos       Position        // Current position
    Inventory []Item          // Player inventory
    Implants  [5]Implant     // Cybernetic implants
    MaxInv    int            // Maximum inventory size
}
```

**Methods:**

##### `func (p *Player) Move(direction rune, width, height int)`
Moves the player in the specified direction within bounds.

**Parameters:**
- `direction`: Movement direction ('↑', '↓', '←', '→')
- `width, height`: Boundary constraints

##### `func (p *Player) AddItemToInventory(item Item) bool`
Adds an item to the player's inventory if space is available.

**Returns:** `true` if item was added successfully

##### `func (p *Player) RemoveItemFromInventory(index int) bool`
Removes an item from inventory by index.

**Returns:** `true` if removal was successful

##### `func (p *Player) GetPosition() (int, int)`
Returns the player's current X, Y coordinates.

#### `type Class`
Defines a character class with base statistics.

```go
type Class struct {
    Name        string  // Class name
    Description string  // Class description
    MaxHP       int     // Base health points
    Force       int     // Base attack power
    Speed       int     // Base speed
    Defense     int     // Base defense
    Accuracy    int     // Base accuracy
}
```

#### `type PlayerStats`
Contains current player statistics and progression.

```go
type PlayerStats struct {
    Level        int     // Current level
    Exp          float32 // Current experience
    NextLevelExp int     // Experience needed for next level
    Force        int     // Current force (attack)
    Speed        int     // Current speed
    Defense      int     // Current defense
    Accuracy     int     // Current accuracy
    MaxHP        int     // Maximum health
    CurrentHP    int     // Current health
}
```

### World Types

#### `type World`
Represents a game world containing multiple stages.

```go
type World struct {
    WorldID        int      // Unique world identifier
    Name           string   // World display name
    Stages         []Stage  // Stages within this world
    ClearingReward int      // Reward for completing world
}
```

**Methods:**

##### `func (w *World) GetStage(stageNb int) *Stage`
Retrieves a stage by its number.

**Returns:** Pointer to the stage, or nil if not found

#### `type Stage`
Represents a single stage within a world.

```go
type Stage struct {
    WorldID        int     // Parent world ID
    StageNb        int     // Stage number within world
    Name           string  // Stage display name
    Enemies        []Enemy // Enemies in this stage
    ClearingReward int     // Reward for completing stage
}
```

**Methods:**

##### `func (s *Stage) AdvanceToNextStage(w *World) *Stage`
Returns the next stage in the world, or nil if this is the last stage.

### Item Types

#### `type Item`
Represents a game item with properties.

```go
type Item struct {
    Type        ItemType  // Item category
    Rarity      Rarity   // Item rarity
    Name        string   // Item name
    Description string   // Item description
}
```

#### `type Weapon`
Represents a weapon with attack patterns.

```go
type Weapon struct {
    KeyName string    // Unique weapon identifier
    Type    int       // Weapon type
    Attacks []Attack  // Available attacks
}
```

#### `type Attack`
Defines a weapon attack with damage and timing.

```go
type Attack struct {
    KeyName  string  // Attack identifier
    KeyDesc  string  // Attack description
    Damage   int     // Base damage
    Duration int     // Attack duration (ms)
    CoolDown int     // Cooldown period (ms)
}
```

## ⚔️ Systems Package

### Combat System

#### `type CombatSystem`
Handles all combat-related calculations and mechanics.

#### `func NewCombatSystem() *CombatSystem`
Creates a new combat system instance.

#### `func (cs *CombatSystem) CalculateDamage(attackerForce, targetDefense int) int`
Calculates damage based on attacker force and target defense.

**Parameters:**
- `attackerForce`: Attacker's force stat
- `targetDefense`: Target's defense stat

**Returns:** Final damage amount (minimum 0)

#### `func (cs *CombatSystem) PlayerAttacksEnemy(player *types.Player, enemy *types.Enemy) int`
Processes a player attack against an enemy.

**Returns:** Damage dealt to the enemy

#### `func (cs *CombatSystem) IsEnemyDefeated(enemy *types.Enemy) bool`
Checks if an enemy has been defeated (HP <= 0).

### Inventory System

#### `type InventorySystem`
Manages inventory operations and item handling.

#### `func NewInventorySystem() *InventorySystem`
Creates a new inventory system instance.

#### `func (is *InventorySystem) AddItem(player *types.Player, item types.Item) bool`
Attempts to add an item to the player's inventory.

**Returns:** `true` if successful

#### `func (is *InventorySystem) IsInventoryFull(player *types.Player) bool`
Checks if the player's inventory is at capacity.

### Movement System

#### `type MovementSystem`
Handles player movement and collision detection.

#### `func NewMovementSystem() *MovementSystem`
Creates a new movement system instance.

#### `func (ms *MovementSystem) MovePlayer(player *types.Player, direction rune, width, height int) bool`
Moves the player if the movement is valid.

**Returns:** `true` if position changed

#### `func (ms *MovementSystem) ValidatePosition(x, y, width, height int) bool`
Validates if a position is within game bounds.

## 🏭 Entities Package

### Factory Functions

#### `func NewPlayer(name string, class types.Class, pos types.Position) *types.Player`
Creates a fully initialized player with the specified parameters.

**Parameters:**
- `name`: Player character name
- `class`: Character class configuration
- `pos`: Starting position

**Returns:** Configured Player instance

**Example:**
```go
position := types.Position{X: 1, Y: 1}
player := entities.NewPlayer("Sam", cybersClass, position)
```

#### `func NewWorld(worldID int) *types.World`
Loads and returns a world by ID, using cached data when available.

**Parameters:**
- `worldID`: Unique world identifier

**Returns:** World instance (may be empty if loading fails)

## 📦 Loaders Package

### Asset Loading

#### `func LoadWorlds() error`
Loads all world data from JSON files and caches it in memory.

**Returns:** Error if loading fails

**Usage:** Call once at application startup

#### `func GetWorld(worldID int) (types.World, bool)`
Retrieves a cached world by ID.

**Parameters:**
- `worldID`: World identifier

**Returns:**
- `types.World`: The world data
- `bool`: Whether the world was found

**Example:**
```go
if err := loaders.LoadWorlds(); err != nil {
    log.Fatal(err)
}

world, exists := loaders.GetWorld(1)
if !exists {
    fmt.Println("World not found")
}
```

## ⚙️ Config Package

### Configuration Access

#### Constants
```go
const (
    DefaultPlayerHealth = 100
    MaxInventorySize    = 10
    BaseExpRequirement  = 100
    ExpGrowthRate       = 1.2
    
    DefaultAnimationSpeed = 200 * time.Millisecond
    TickDuration          = 16 * time.Millisecond
    
    DefaultTerminalWidth  = 80
    DefaultTerminalHeight = 24
    HUDHeight             = 5
)
```

#### Default Classes
```go
var DefaultClasses = map[string]ClassConfig{
    "D0C": {
        Name: "D0C",
        Description: "Un robot intelligent, précis et polyvalent.",
        MaxHP: 90, Force: 10, Speed: 12, Defense: 10, Accuracy: 22,
    },
    "APP": {
        Name: "APP",
        Description: "Un robot furtif, rapide et précis.", 
        MaxHP: 80, Force: 14, Speed: 22, Defense: 8, Accuracy: 18,
    },
    "CYBER_SAMURAI": {
        Name: "Cyber-Samurai",
        Description: "A swift and deadly warrior...",
        MaxHP: 100, Force: 15, Speed: 12, Defense: 8, Accuracy: 15,
    },
}
```

#### Asset Paths
```go
var AssetPathsConfig = AssetPaths{
    Root:          "assets",
    DataDir:       "assets/data",
    AnimationsDir: "assets/animations",
    LevelsDir:     "assets/levels",
    // ... more paths
}
```

## 🔧 Usage Examples

### Complete Game Setup
```go
// Initialize game systems
if err := loaders.LoadWorlds(); err != nil {
    log.Fatal("Failed to load worlds:", err)
}

// Create player with selected class
class := config.DefaultClasses["CYBER_SAMURAI"]
gameClass := types.Class{
    Name: class.Name,
    MaxHP: class.MaxHP,
    Force: class.Force,
    // ... other fields
}

// Start new game
game := NewGameInstance(gameClass)

// Game loop
for {
    // Handle input, update systems, render
    location, worldID := game.CurrentLocation()
    fmt.Printf("Location: %s (World %d)\n", location, worldID)
    
    // Advance if conditions met
    if shouldAdvance {
        success := game.Advance()
        if !success {
            fmt.Println("Cannot advance further")
        }
    }
}
```

### Combat Example
```go
// Combat scenario
combat := systems.NewCombatSystem()
damage := combat.PlayerAttacksEnemy(game.Player, enemy)
fmt.Printf("Player deals %d damage!\n", damage)

if combat.IsEnemyDefeated(enemy) {
    fmt.Println("Enemy defeated!")
    // Add experience, loot, etc.
}
```

### Inventory Management
```go
// Add item to inventory
inventory := systems.NewInventorySystem()
item := types.Item{
    Name: "Health Potion",
    Type: types.Consumable,
    // ... other fields
}

if inventory.AddItem(game.Player, item) {
    fmt.Println("Item added to inventory")
} else {
    fmt.Println("Inventory full!")
}
```

This API provides a clean, well-documented interface for all game functionality while maintaining the modular architecture.