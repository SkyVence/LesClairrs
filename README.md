# 🎮 ProjectRed RPG

A cyberpunk-themed terminal-based RPG developed in Go, featuring a modern architecture with clean separation of concerns.

## 🏗️ Project Structure

```
src/
├── main.go                    # Entry point
├── go.mod & go.sum           # Dependencies
│
├── config/                   # Configuration management
│   ├── constants.go          # Game balance, defaults, class configs
│   ├── paths.go             # Asset paths configuration
│   └── engine.go            # Engine configuration
│
├── engine/                   # Game engine core
│   ├── animation.go         # Animation system
│   ├── engine.go           # Main engine wrapper
│   ├── input.go            # Input handling
│   ├── program.go          # Program lifecycle
│   ├── renderer.go         # Rendering system
│   └── tea.go              # Terminal UI framework
│
├── game/                     # Game logic
│   ├── game.go              # Main game coordinator
│   ├── render.go            # Game rendering logic
│   │
│   ├── types/               # Type definitions
│   │   ├── player.go        # Player, Class, PlayerStats + methods
│   │   ├── world.go         # World, Stage, Position + methods  
│   │   ├── enemy.go         # Enemy types
│   │   ├── items.go         # Item, Weapon, Attack types
│   │   └── enums.go         # Enums (BodyParts, Rarity, ItemType)
│   │
│   ├── entities/            # Entity creation and management
│   │   ├── player.go        # Player factory functions
│   │   └── world.go         # World factory functions
│   │
│   ├── systems/             # Game systems
│   │   ├── combat.go        # Combat calculations and logic
│   │   ├── inventory.go     # Inventory management
│   │   └── movement.go      # Movement and collision detection
│   │
│   └── loaders/             # Asset loading
│       └── loadLevels.go    # World/level data loading
│
├── ui/                       # User interface components
│   ├── hud.go              # Heads-up display
│   ├── menu.go             # Menu systems
│   └── spinner.go          # Loading indicators
│
└── assets/                   # Game assets
    ├── data/               # Game data (weapons, enemies, etc.)
    ├── animations/         # Animation files
    ├── interface/          # UI localization files
    └── levels/             # World and level definitions
```

## 🚀 Installation

```bash
git clone https://github.com/SkyVence/projet-red_rpg
cd projet-red_rpg
go mod tidy
go run src/main.go
```

## 🎮 Controls

- `↑`/`↓` `←`/`→` : Navigation
- `Enter` : Select
- `Esc` : Back/Quit

## 🧩 Architecture Overview

### **Types Package** (`game/types/`)
Contains all core type definitions with their associated methods:
- **Player**: Character data, stats, inventory management
- **World/Stage**: Level structure and progression
- **Enemy**: Opponent definitions and behavior
- **Items**: Equipment and consumables

### **Systems Package** (`game/systems/`)
Implements game logic using a systems architecture:
- **CombatSystem**: Damage calculations, battle mechanics
- **InventorySystem**: Item management, equipment
- **MovementSystem**: Player movement, collision detection

### **Entities Package** (`game/entities/`)
Factory functions for creating and initializing game entities:
- Player creation with class selection
- World loading and initialization

### **Config Package** (`config/`)
Centralized configuration management:
- Game balance constants
- Asset paths
- Engine settings

## 🎯 Key Features

- **Clean Architecture**: Separation of concerns with dedicated packages
- **Systems-Based Design**: Modular game logic systems
- **Type Safety**: Strong typing with Go's type system
- **Asset Management**: Structured asset loading and caching
- **Terminal UI**: Modern terminal interface with animations
- **Configurable**: Easy-to-modify game balance and settings

## 🛠️ Development

### Building
```bash
cd src
go build .
```

### Running Tests
```bash
go test ./...
```

### Adding New Features

1. **New Game Systems**: Add to `game/systems/`
2. **New Entity Types**: Define in `game/types/`, create in `game/entities/`
3. **New Assets**: Place in appropriate `assets/` subdirectory
4. **Configuration**: Update `config/` files for new settings

## 📚 Documentation

- [Architecture Guide](docs/architecture.md) - Detailed architecture explanation
- [API Reference](docs/api.md) - Package and function documentation
- [Asset Format](docs/assets.md) - Asset file format specifications

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the existing architecture patterns
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is open source. See LICENSE file for details.

---

**Développé en Go** | **Interface Terminal Moderne** | **Architecture Modulaire**
