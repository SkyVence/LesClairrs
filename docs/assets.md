# 📁 Asset Format Specification

This document describes the file formats and structure used for game assets in ProjectRed RPG.

## 📋 Table of Contents

- [Asset Directory Structure](#asset-directory-structure)
- [World Data Format](#world-data-format)
- [Animation Format](#animation-format)
- [Weapon Data Format](#weapon-data-format)
- [Localization Format](#localization-format)
- [Best Practices](#best-practices)

## 📂 Asset Directory Structure

```
assets/
├── data/                    # Game data files
│   ├── iron-sword.json     # Weapon definitions
│   ├── steel-dagger.json
│   └── wooden-bow.json
├── animations/              # Animation files
│   ├── empty.anim          # Empty/idle animation
│   ├── loader.anim         # Loading spinner
│   └── player-running.anim # Player movement
├── interface/               # UI localization
│   └── fr.json             # French text
├── levels/                  # World and level data
│   ├── world-1.json        # First world definition
│   └── world-2.json        # Second world definition
└── logo.txt                # ASCII art logo
```

## 🌍 World Data Format

World files define game levels, stages, and enemies. They use JSON format for easy editing and parsing.

### File Location
`assets/levels/world-{id}.json`

### Schema
```json
{
  "WorldID": 1,
  "Name": "World Display Name",
  "Stages": [
    {
      "StageNb": 1,
      "Name": "Stage Display Name", 
      "Enemies": [
        {
          "Name": "Enemy Name",
          "Force": 5,
          "Speed": 5,
          "Defense": 3,
          "Accuracy": 7,
          "MaxHP": 20,
          "CurrentHP": 20,
          "ExpReward": 20
        }
      ],
      "ClearingReward": 50
    }
  ],
  "ClearingReward": 100
}
```

### Example: `world-1.json`
```json
{
  "WorldID": 1,
  "Name": "Level 1 - Testing Grounds",
  "Stages": [
    {
      "StageNb": 1,
      "Name": "Stage 1 - The Beginning",
      "Enemies": [
        {
          "Name": "Rogue Drone",
          "Force": 5,
          "Speed": 5,
          "Defense": 3,
          "Accuracy": 7,
          "MaxHP": 20,
          "CurrentHP": 20,
          "ExpReward": 20
        },
        {
          "Name": "Street Thug",
          "Force": 6,
          "Speed": 4,
          "Defense": 4,
          "Accuracy": 6,
          "MaxHP": 25,
          "CurrentHP": 25,
          "ExpReward": 25
        }
      ],
      "ClearingReward": 50
    },
    {
      "StageNb": 2,
      "Name": "Stage 2 - Deeper into the City",
      "Enemies": [
        {
          "Name": "Cyber Hound",
          "Force": 8,
          "Speed": 7,
          "Defense": 5,
          "Accuracy": 8,
          "MaxHP": 30,
          "CurrentHP": 30,
          "ExpReward": 40
        }
      ],
      "ClearingReward": 75
    }
  ],
  "ClearingReward": 200
}
```

### Field Descriptions

#### World Level
- `WorldID`: Unique integer identifier for the world
- `Name`: Display name shown to players
- `Stages`: Array of stage definitions within this world
- `ClearingReward`: Experience/currency awarded for completing entire world

#### Stage Level
- `StageNb`: Stage number within the world (starts at 1)
- `Name`: Display name for the stage
- `Enemies`: Array of enemy definitions for this stage
- `ClearingReward`: Reward for completing this specific stage

#### Enemy Level
- `Name`: Enemy display name
- `Force`: Attack power (used in damage calculations)
- `Speed`: Movement/initiative speed
- `Defense`: Damage reduction capability
- `Accuracy`: Hit chance modifier
- `MaxHP`/`CurrentHP`: Health points (usually equal at start)
- `ExpReward`: Experience points awarded when defeated

## 🎬 Animation Format

Animation files use a simple text format with frame separation.

### File Location
`assets/animations/{name}.anim`

### Format
- Each frame is separated by `---` on its own line
- Frames can contain any ASCII characters
- Leading/trailing newlines are trimmed
- Spaces within frames are preserved

### Example: `player-running.anim`
```
  O
 /|\
 / \
---
  O
 <|/
 / >
```

### Loading in Code
```go
frames, err := engine.LoadAnimationFile("player-running.anim")
if err != nil {
    log.Fatal("Failed to load animation:", err)
}

animation := engine.NewAnimation(frames)
```

## ⚔️ Weapon Data Format

Weapon files define combat equipment with attack patterns.

### File Location
`assets/data/{weapon-name}.json`

### Schema
```json
{
  "KeyName": "weapon-identifier",
  "Type": 0,
  "Attacks": [
    {
      "KeyName": "attack-id",
      "KeyDesc": "Attack description",
      "Damage": 15,
      "Duration": 500,
      "CoolDown": 1000
    }
  ]
}
```

### Example: `iron-sword.json`
```json
{
  "KeyName": "iron-sword",
  "Type": 0,
  "Attacks": [
    {
      "KeyName": "slash",
      "KeyDesc": "A quick sword slash",
      "Damage": 15,
      "Duration": 500,
      "CoolDown": 1000
    },
    {
      "KeyName": "thrust", 
      "KeyDesc": "A powerful thrust attack",
      "Damage": 25,
      "Duration": 800,
      "CoolDown": 2000
    }
  ]
}
```

### Field Descriptions
- `KeyName`: Unique identifier for the weapon
- `Type`: Weapon category (0 = melee, 1 = ranged, etc.)
- `Attacks`: Array of available attack patterns
  - `KeyName`: Unique attack identifier
  - `KeyDesc`: Human-readable attack description
  - `Damage`: Base damage value
  - `Duration`: Animation/execution time in milliseconds
  - `CoolDown`: Minimum time between uses in milliseconds

## 🌐 Localization Format

Localization files provide translated text for the user interface.

### File Location
`assets/interface/{language}.json`

### Schema
```json
{
  "ui": {
    "buttons": {
      "start": "Translated Start",
      "quit": "Translated Quit"
    },
    "messages": {
      "welcome": "Translated welcome message"
    }
  },
  "game": {
    "actions": {
      "move": "Translated move command"
    }
  }
}
```

### Example: `fr.json`
```json
{
  "ui": {
    "buttons": {
      "start": "Commencer",
      "settings": "Paramètres", 
      "quit": "Quitter"
    },
    "messages": {
      "welcome": "Bienvenue dans ProjectRed RPG!",
      "loading": "Chargement..."
    }
  },
  "game": {
    "actions": {
      "move": "Déplacer",
      "attack": "Attaquer",
      "inventory": "Inventaire"
    },
    "status": {
      "health": "Santé",
      "level": "Niveau",
      "experience": "Expérience"
    }
  }
}
```

## 🎨 ASCII Art Format

ASCII art files are plain text files containing stylized text graphics.

### File Location
`assets/logo.txt` or similar

### Format
- Plain text files with ASCII characters
- No special formatting required
- Used for logos, banners, decorative elements

### Example: `logo.txt`
```
 ██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗
 ██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝
 ██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║   
 ██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║   
 ██║     ██║  ██║╚██████╔╝╚█████╔╝███████╗╚██████╗   ██║   
 ╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝   
                                                            
             ██████╗ ███████╗██████╗     ██████╗ ██████╗  ██████╗ 
             ██╔══██╗██╔════╝██╔══██╗    ██╔══██╗██╔══██╗██╔════╝ 
             ██████╔╝█████╗  ██║  ██║    ██████╔╝██████╔╝██║  ███╗
             ██╔══██╗██╔══╝  ██║  ██║    ██╔══██╗██╔═══╝ ██║   ██║
             ██║  ██║███████╗██████╔╝    ██║  ██║██║     ╚██████╔╝
             ╚═╝  ╚═╝╚══════╝╚═════╝     ╚═╝  ╚═╝╚═╝      ╚═════╝ 
```

## ✅ Best Practices

### File Naming
- Use lowercase with hyphens: `iron-sword.json`
- Include clear identifiers: `world-1.json`
- Use descriptive names: `player-running.anim`

### JSON Formatting
- Use consistent indentation (2 or 4 spaces)
- Include all required fields
- Validate JSON syntax before committing
- Use meaningful field names

### Content Guidelines
- Keep file sizes reasonable (< 100KB for data files)
- Use consistent naming conventions within files
- Include comments in documentation, not in JSON files
- Test all assets in-game before finalizing

### Version Control
- Commit asset files with descriptive messages
- Group related changes (e.g., "Add level 3 enemies")
- Avoid committing binary assets when possible
- Keep ASCII art in separate files for easier editing

### Performance Considerations
- Load assets at startup when possible
- Cache parsed data in memory
- Use efficient data structures for lookups
- Minimize file I/O during gameplay

### Validation
Always validate asset files:
```bash
# JSON validation
jq . assets/levels/world-1.json

# Check for required fields
grep -q "WorldID" assets/levels/world-1.json
```

### Adding New Asset Types
1. Define the data structure in `game/types/`
2. Create loading logic in `game/loaders/`
3. Add file paths to `config/paths.go`
4. Document the format in this file
5. Create example assets
6. Add validation tests

This asset system provides flexibility while maintaining structure and consistency across all game content.