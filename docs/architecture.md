# 🏗️ Architecture Guide

This document explains the architecture and design patterns used in ProjectRed RPG.

## 📋 Table of Contents

- [Overview](#overview)
- [Package Structure](#package-structure)
- [Design Patterns](#design-patterns)
- [Data Flow](#data-flow)
- [Dependencies](#dependencies)
- [Best Practices](#best-practices)

## 🎯 Overview

ProjectRed RPG follows a **modular, systems-based architecture** that promotes:

- **Separation of Concerns**: Each package has a single, well-defined responsibility
- **Loose Coupling**: Minimal dependencies between packages
- **High Cohesion**: Related functionality is grouped together
- **Scalability**: Easy to add new features without restructuring
- **Testability**: Each component can be tested independently

## 📦 Package Structure

### **Core Packages**

#### `config/` - Configuration Management
**Purpose**: Centralized configuration and constants
**Contents**:
- `constants.go`: Game balance, timing, display settings
- `paths.go`: Asset file paths and directory structure
- `engine.go`: Engine-specific configuration

**Design Pattern**: Singleton pattern with global configuration instances

#### `engine/` - Game Engine
**Purpose**: Low-level engine functionality and terminal management
**Contents**:
- Terminal I/O handling
- Rendering pipeline
- Input processing
- Animation system
- Program lifecycle

**Design Pattern**: Facade pattern hiding engine complexity

#### `game/types/` - Type Definitions
**Purpose**: Core data structures and their methods
**Contents**:
- Player, Enemy, World, Stage definitions
- Item and equipment types
- Enums and constants

**Design Pattern**: Domain-driven design with rich domain objects

#### `game/entities/` - Entity Factories
**Purpose**: Object creation and initialization
**Contents**:
- Player creation with class setup
- World loading and initialization
- Entity lifecycle management

**Design Pattern**: Factory pattern for object creation

#### `game/systems/` - Game Logic Systems
**Purpose**: Modular game logic implementation
**Contents**:
- Combat calculations and mechanics
- Inventory management operations
- Movement and collision detection

**Design Pattern**: Entity-Component-System (ECS) inspired architecture

#### `game/loaders/` - Asset Management
**Purpose**: Loading and caching game data
**Contents**:
- JSON file parsing
- Asset caching
- Data validation

**Design Pattern**: Repository pattern with caching

#### `ui/` - User Interface
**Purpose**: Terminal UI components and rendering
**Contents**:
- HUD display
- Menu systems
- Loading indicators

**Design Pattern**: Component-based UI architecture

## 🔄 Design Patterns

### 1. **Systems Architecture**
Game logic is split into independent systems that operate on shared data:

```go
type Game struct {
    // Data
    Player       *types.Player
    CurrentWorld *types.World
    
    // Systems
    Combat    *systems.CombatSystem
    Inventory *systems.InventorySystem
    Movement  *systems.MovementSystem
}
```

**Benefits**:
- Modular logic that's easy to test
- Systems can be enabled/disabled independently
- Clear separation between data and behavior

### 2. **Factory Pattern**
Entity creation is handled by dedicated factory functions:

```go
// entities/player.go
func NewPlayer(name string, class types.Class, pos types.Position) *types.Player

// entities/world.go  
func NewWorld(worldID int) *types.World
```

**Benefits**:
- Centralized object creation logic
- Consistent initialization
- Easy to modify creation process

### 3. **Repository Pattern**
Asset loading uses a repository pattern with caching:

```go
// loaders/loadLevels.go
func LoadWorlds() error           // Load data
func GetWorld(id int) (World, bool) // Retrieve cached data
```

**Benefits**:
- Abstracted data access
- Built-in caching
- Easy to swap data sources

### 4. **Configuration Singleton**
Global configuration accessible throughout the application:

```go
// config/constants.go
var DefaultClasses = map[string]ClassConfig{...}

// config/paths.go
var AssetPathsConfig = DefaultAssetPaths()
```

**Benefits**:
- Centralized configuration
- Easy to modify settings
- Consistent access pattern

## 🌊 Data Flow

### **Game Initialization**
1. `main.go` → Creates game instance
2. `game/render.go` → Initializes UI components
3. `game/entities/` → Creates player and world
4. `game/loaders/` → Loads asset data
5. `game/systems/` → Initializes game systems

### **Game Loop**
1. **Input** → `engine/input.go` processes user input
2. **Update** → Game systems process the input:
   - `systems/movement.go` handles player movement
   - `systems/combat.go` processes combat
   - `systems/inventory.go` manages items
3. **Render** → `game/render.go` creates display output
4. **Display** → `engine/renderer.go` outputs to terminal

### **Asset Loading**
1. `config/paths.go` → Provides asset file paths
2. `loaders/loadLevels.go` → Reads JSON files
3. Cache in memory for fast access
4. `entities/world.go` → Uses cached data

## 🔗 Dependencies

### **Dependency Graph**
```
main.go
├── engine/ (low-level terminal handling)
├── game/ (game logic coordinator)
│   ├── types/ (data structures)
│   ├── entities/ (object creation)
│   │   └── → types/, loaders/
│   ├── systems/ (game logic)
│   │   └── → types/
│   └── loaders/ (asset management)
│       └── → types/
├── ui/ (interface components)
│   └── → engine/, game/types/
└── config/ (configuration, no dependencies)
```

### **Import Rules**
- `config/` has no dependencies (base layer)
- `types/` only imports `config/` 
- `entities/` can import `types/`, `loaders/`
- `systems/` only import `types/`
- `game/` coordinates all packages
- Circular dependencies are avoided

## ✅ Best Practices

### **Adding New Features**

#### New Game System
1. Create file in `game/systems/`
2. Define system struct with `New*System()` constructor
3. Add methods that operate on `types.*` structures
4. Integrate into `game/game.go`

```go
// game/systems/magic.go
type MagicSystem struct{}

func NewMagicSystem() *MagicSystem { return &MagicSystem{} }

func (ms *MagicSystem) CastSpell(player *types.Player, spell types.Spell) error {
    // Implementation
}
```

#### New Entity Type
1. Define types in `game/types/`
2. Add methods to the type
3. Create factory in `game/entities/`
4. Update systems as needed

```go
// game/types/npc.go
type NPC struct {
    Name     string
    Dialogue []string
}

func (npc *NPC) Speak() string { /* implementation */ }

// game/entities/npc.go
func NewNPC(name string, dialogue []string) *types.NPC {
    return &types.NPC{Name: name, Dialogue: dialogue}
}
```

#### New Configuration
1. Add constants to `config/constants.go`
2. Add paths to `config/paths.go` if needed
3. Use throughout the application

### **Testing Strategy**
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test system interactions
- **Package Tests**: Test package boundaries

```go
// game/systems/combat_test.go
func TestCombatSystem_CalculateDamage(t *testing.T) {
    cs := NewCombatSystem()
    damage := cs.CalculateDamage(10, 5)
    assert.Equal(t, 5, damage)
}
```

### **Error Handling**
- Use Go's error return pattern
- Wrap errors with context
- Handle errors at appropriate levels

```go
func LoadWorlds() error {
    if err := loadFromFile(); err != nil {
        return fmt.Errorf("failed to load worlds: %w", err)
    }
    return nil
}
```

## 🚀 Performance Considerations

- **Asset Caching**: All game data is loaded once and cached
- **Memory Management**: Reuse objects where possible
- **Rendering Optimization**: Only render when state changes
- **Input Debouncing**: Prevent rapid repeated inputs

## 🔮 Future Enhancements

The architecture supports easy addition of:
- **Save/Load System**: Add to `loaders/` package
- **Multiplayer**: Extend systems to handle multiple players
- **Scripting**: Add script execution system
- **Mod Support**: Plugin architecture for systems
- **Audio**: Add audio system to engine
- **Graphics**: Enhanced rendering with color/effects

This architecture provides a solid foundation for a scalable, maintainable game engine while keeping the codebase clean and organized.