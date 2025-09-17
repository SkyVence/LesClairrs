# Project Red Engine API Documentation

A comprehensive guide to the terminal-based game engine for Project Red RPG.

## Table of Contents

1. [Overview](#overview)
2. [Core Engine Architecture](#core-engine-architecture)
3. [Program Management](#program-management)
4. [Rendering System](#rendering-system)
5. [Input Handling](#input-handling)
6. [Animation System](#animation-system)
7. [Internationalization](#internationalization)
8. [Complete Examples](#complete-examples)
9. [Best Practices](#best-practices)

## Overview

The Project Red Engine is a terminal-based game engine built in Go, designed around the Model-View-Update (MVU) architecture. It provides a complete framework for creating interactive terminal applications with rendering, input handling, animations, and internationalization support.

### Key Features

- **Terminal-based rendering** with advanced terminal control
- **Real-time input handling** including keyboard and special keys
- **Animation system** with frame-based animations
- **Internationalization support** with JSON-based language files
- **Coordinate-based positioning** for precise pixel/character placement
- **Alternate screen support** for full-screen applications

### Architecture

The engine follows the MVU pattern where:
- **Model**: Application state
- **View**: String representation of the current state
- **Update**: State transitions based on messages

## Core Engine Architecture

### Game Interface

All games must implement the `Game` interface:

```go
type Game interface {
    Init() Msg
    Update(Msg) (Model, Cmd)
    View() string
}
```

**Implementation Example:**

```go
package mygame

import "projectred-rpg.com/engine"

type MyGame struct {
    player   Player
    score    int
    width    int
    height   int
}

func (g *MyGame) Init() engine.Msg {
    // Initialize game state
    return nil
}

func (g *MyGame) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg := msg.(type) {
    case engine.KeyMsg:
        // Handle key input
        switch msg.Rune {
        case 'q':
            return g, engine.Quit
        case '↑':
            g.player.MoveUp()
        }
    case engine.SizeMsg:
        g.width = msg.Width
        g.height = msg.Height
    }
    return g, nil
}

func (g *MyGame) View() string {
    // Return string representation of game state
    return g.renderGameWorld()
}
```

### Engine Wrapper

Use `engine.Wrap()` to convert your game into an engine-compatible model:

```go
game := &MyGame{}
model := engine.Wrap(game)
```

## Program Management

### Creating and Running Programs

The `Program` struct manages the application lifecycle:

```go
func NewProgram(model Model, opts ...ProgramOption) *Program
```

**Basic Usage:**

```go
package main

import (
    "log"
    "projectred-rpg.com/engine"
)

func main() {
    game := NewMyGame()
    program := engine.NewProgram(
        engine.Wrap(game),
        engine.WithAltScreen(), // Optional: use alternate screen
    )
    
    if err := program.Run(); err != nil {
        log.Fatalf("Program error: %v", err)
    }
}
```

### Program Options

- `WithAltScreen()`: Enables alternate screen buffer for full-screen applications

### Terminal Size Detection

Get terminal dimensions:

```go
width, height := program.GetSize()
```

The engine automatically sends `SizeMsg` when terminal is resized.

## Rendering System

### Coordinate System

The engine uses a character-based coordinate system where:
- Origin (0,0) is at the top-left corner
- X increases to the right
- Y increases downward
- Each position represents one terminal character

### Renderer Interface

The renderer provides advanced terminal control:

```go
type Renderer interface {
    Write(string)                    // Write content to screen
    ClearScreen()                   // Clear entire screen
    SetCursor(x, y int)            // Position cursor at coordinates
    ShowCursor() / HideCursor()    // Control cursor visibility
    EnterAltScreen() / ExitAltScreen() // Alternate screen control
    SetWindowTitle(string)         // Set terminal window title
}
```

### Pixel Placement in Game Space

To place characters at specific coordinates:

```go
func PlacePixelAt(x, y int, char rune, grid [][]rune, maxWidth, maxHeight int) {
    if x >= 0 && x < maxWidth && y >= 0 && y < maxHeight {
        grid[y][x] = char
    }
}
```

**Example: Drawing a Border**

```go
func DrawBorder(grid [][]rune, width, height int) {
    for i := 0; i < height; i++ {
        for j := 0; j < width; j++ {
            if i == 0 || i == height-1 || j == 0 || j == width-1 {
                var char rune
                switch {
                case i == 0 && j == 0:
                    char = '┌'
                case i == 0 && j == width-1:
                    char = '┐'
                case i == height-1 && j == 0:
                    char = '└'
                case i == height-1 && j == width-1:
                    char = '┘'
                case i == 0 || i == height-1:
                    char = '─'
                default:
                    char = '│'
                }
                grid[i][j] = char
            }
        }
    }
}
```

**Example: Sprite Rendering**

```go
func RenderSprite(sprite string, x, y int, grid [][]rune, maxWidth, maxHeight int) {
    lines := strings.Split(sprite, "\n")
    for i, line := range lines {
        for j, char := range line {
            spriteX := x + j
            spriteY := y + i
            if spriteX >= 0 && spriteX < maxWidth && spriteY >= 0 && spriteY < maxHeight {
                grid[spriteY][spriteX] = char
            }
        }
    }
}
```

### Grid to String Conversion

Convert your 2D grid to a string for rendering:

```go
func GridToString(grid [][]rune) string {
    var builder strings.Builder
    for _, row := range grid {
        builder.WriteString(string(row))
        builder.WriteString("\n")
    }
    return strings.TrimRight(builder.String(), "\n")
}
```

## Input Handling

### Message Types

The engine provides several message types:

```go
type KeyMsg struct {
    Rune rune  // The character pressed
}

type QuitMsg struct{}     // Application quit signal

type SizeMsg struct {     // Terminal resize
    Width  int
    Height int
}

type TickMsg struct {     // Timer tick
    Time time.Time
}
```

### Key Input Handling

**Basic Key Detection:**

```go
func (g *MyGame) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg := msg.(type) {
    case engine.KeyMsg:
        switch msg.Rune {
        case 'q', 'Q':
            return g, engine.Quit()
        case 'w', 'W', '↑':
            g.player.MoveUp()
        case 's', 'S', '↓':
            g.player.MoveDown()
        case 'a', 'A', '←':
            g.player.MoveLeft()
        case 'd', 'D', '→':
            g.player.MoveRight()
        case ' ':  // Spacebar
            g.player.Action()
        case '\r', '\n':  // Enter key
            g.confirm()
        case 3:  // Ctrl+C
            return g, engine.Quit()
        }
    }
    return g, nil
}
```

### Special Keys

The engine automatically converts escape sequences to Unicode arrows:
- `↑` (Up arrow)
- `↓` (Down arrow) 
- `←` (Left arrow)
- `→` (Right arrow)

### Player Movement Example

```go
type Position struct {
    X, Y int
}

func (p *Player) Move(direction rune, maxWidth, maxHeight int) {
    switch direction {
    case '↑', 'w', 'W':
        if p.Position.Y > 1 {  // Respect border
            p.Position.Y--
        }
    case '↓', 's', 'S':
        if p.Position.Y < maxHeight-2 {  // Respect border
            p.Position.Y++
        }
    case '←', 'a', 'A':
        if p.Position.X > 1 {  // Respect border
            p.Position.X--
        }
    case '→', 'd', 'D':
        if p.Position.X < maxWidth-2 {  // Respect border
            p.Position.X++
        }
    }
}
```

## Animation System

### Loading Animation Files

Animation files use `---` as frame separators:

```
Frame 1 content
here
---
Frame 2 content
here
---
Frame 3 content
here
```

**Loading Function:**

```go
frames, err := engine.LoadAnimationFile("assets/animations/player-running.anim")
if err != nil {
    log.Printf("Failed to load animation: %v", err)
    return
}
```

### Animation Structure

```go
type Animation struct {
    Frames []string
    // private fields for state management
}

func NewAnimation(frames []string) Animation
func (a Animation) Init() Cmd
func (a Animation) Update(msg Msg) (Animation, Cmd)
func (a Animation) View() string
```

### Using Animations

**Complete Animation Example:**

```go
type GameWithAnimation struct {
    playerAnim engine.Animation
    position   Position
}

func (g *GameWithAnimation) Init() engine.Msg {
    // Load animation frames
    frames, err := engine.LoadAnimationFile("assets/animations/player.anim")
    if err != nil {
        frames = []string{"@"}  // Fallback
    }
    
    g.playerAnim = engine.NewAnimation(frames)
    return g.playerAnim.Init()  // Start animation
}

func (g *GameWithAnimation) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    // Update animation
    var cmd engine.Cmd
    g.playerAnim, cmd = g.playerAnim.Update(msg)
    
    switch msg := msg.(type) {
    case engine.KeyMsg:
        // Handle other input...
    }
    
    return g, cmd
}

func (g *GameWithAnimation) View() string {
    // Get current animation frame
    sprite := g.playerAnim.View()
    
    // Render sprite at position
    grid := make([][]rune, height)
    // ... initialize grid ...
    
    RenderSprite(sprite, g.position.X, g.position.Y, grid, width, height)
    return GridToString(grid)
}
```

### Custom Animation Timing

Animations default to 200ms per frame. For custom timing, modify the Animation struct or create a wrapper.

## Internationalization

### Language File Format

JSON files in `assets/interface/`:

```json
{
    "game.title": "Project Red RPG",
    "player.health": "Health: {current}/{max}",
    "level.world1.name": "Forest of Beginnings",
    "level.world1.stage1": "Misty Path",
    "menu.start": "Start Game",
    "menu.quit": "Quit"
}
```

### Loading and Using Languages

```go
// Load language catalog
lang, err := engine.Load("fr")  // loads assets/interface/fr.json
if err != nil {
    log.Printf("Failed to load language: %v", err)
    return
}

// Simple text retrieval
title := lang.Text("game.title")  // "Project Red RPG"

// Text with placeholders
health := lang.Text("player.health", 75, 100)  // "Health: 75/100"

// Missing keys return bracketed key
missing := lang.Text("nonexistent.key")  // "⟦nonexistent.key⟧"
```

### Placeholder System

Placeholders are replaced in encounter order:

```go
// JSON: "combat.damage": "Player deals {damage} damage to {enemy}!"
text := lang.Text("combat.damage", 25, "Goblin")
// Result: "Player deals 25 damage to Goblin!"
```

### Dynamic Language Content

```go
// Example: Dynamic world/stage names
func GetLocationName(lang engine.Catalog, worldID, stageID int) string {
    worldKey := fmt.Sprintf("level.world%d.name", worldID)
    stageKey := fmt.Sprintf("level.world%d.stage%d", worldID, stageID)
    
    worldName := lang.Text(worldKey)
    stageName := lang.Text(stageKey)
    
    return fmt.Sprintf("%s - %s", worldName, stageName)
}
```

## Complete Examples

### Simple Game Template

```go
package main

import (
    "log"
    "strings"
    "projectred-rpg.com/engine"
)

type SimpleGame struct {
    playerX, playerY int
    width, height    int
}

func NewSimpleGame() *SimpleGame {
    return &SimpleGame{
        playerX: 5,
        playerY: 5,
    }
}

func (g *SimpleGame) Init() engine.Msg {
    return nil
}

func (g *SimpleGame) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg := msg.(type) {
    case engine.SizeMsg:
        g.width = msg.Width
        g.height = msg.Height
    case engine.KeyMsg:
        switch msg.Rune {
        case 'q':
            return g, engine.Quit()
        case '↑':
            if g.playerY > 1 {
                g.playerY--
            }
        case '↓':
            if g.playerY < g.height-2 {
                g.playerY++
            }
        case '←':
            if g.playerX > 1 {
                g.playerX--
            }
        case '→':
            if g.playerX < g.width-2 {
                g.playerX++
            }
        }
    }
    return g, nil
}

func (g *SimpleGame) View() string {
    if g.width <= 0 || g.height <= 0 {
        return "Initializing..."
    }
    
    // Create grid
    grid := make([][]rune, g.height)
    for i := range grid {
        grid[i] = make([]rune, g.width)
        for j := range grid[i] {
            grid[i][j] = ' '
        }
    }
    
    // Draw border
    DrawBorder(grid, g.width, g.height)
    
    // Draw player
    grid[g.playerY][g.playerX] = '@'
    
    // Convert to string
    return GridToString(grid)
}

func main() {
    game := NewSimpleGame()
    program := engine.NewProgram(
        engine.Wrap(game),
        engine.WithAltScreen(),
    )
    
    if err := program.Run(); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
```

### Advanced Game with Animation and Language

```go
type AdvancedGame struct {
    playerX, playerY int
    width, height    int
    playerAnim       engine.Animation
    lang            engine.Catalog
}

func NewAdvancedGame() *AdvancedGame {
    // Load language
    lang, err := engine.Load("en")
    if err != nil {
        lang = make(engine.Catalog)  // Empty fallback
    }
    
    // Load animation
    frames, err := engine.LoadAnimationFile("assets/animations/player.anim")
    if err != nil {
        frames = []string{"@"}  // Fallback
    }
    
    return &AdvancedGame{
        playerX:    10,
        playerY:    10,
        lang:      lang,
        playerAnim: engine.NewAnimation(frames),
    }
}

func (g *AdvancedGame) Init() engine.Msg {
    return g.playerAnim.Init()
}

func (g *AdvancedGame) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    // Update animation
    var cmd engine.Cmd
    g.playerAnim, cmd = g.playerAnim.Update(msg)
    
    switch msg := msg.(type) {
    case engine.SizeMsg:
        g.width = msg.Width
        g.height = msg.Height
    case engine.KeyMsg:
        switch msg.Rune {
        case 'q':
            return g, engine.Quit()
        case '↑', '↓', '←', '→':
            g.movePlayer(msg.Rune)
        }
    }
    return g, cmd
}

func (g *AdvancedGame) movePlayer(direction rune) {
    switch direction {
    case '↑':
        if g.playerY > 1 {
            g.playerY--
        }
    case '↓':
        if g.playerY < g.height-2 {
            g.playerY++
        }
    case '←':
        if g.playerX > 1 {
            g.playerX--
        }
    case '→':
        if g.playerX < g.width-2 {
            g.playerX++
        }
    }
}

func (g *AdvancedGame) View() string {
    if g.width <= 0 || g.height <= 0 {
        return g.lang.Text("loading", "Initializing...")
    }
    
    // Create grid
    grid := make([][]rune, g.height)
    for i := range grid {
        grid[i] = make([]rune, g.width)
        for j := range grid[i] {
            grid[i][j] = ' '
        }
    }
    
    // Draw border
    DrawBorder(grid, g.width, g.height)
    
    // Draw animated player sprite
    sprite := g.playerAnim.View()
    RenderSprite(sprite, g.playerX, g.playerY, grid, g.width, g.height)
    
    return GridToString(grid)
}
```

## Best Practices

### 1. Coordinate Management

- Always validate coordinates before placing characters
- Respect border boundaries in movement logic
- Use consistent coordinate systems across your application

### 2. Performance Optimization

- Minimize string allocations in the View() method
- Reuse grid buffers when possible
- Cache language strings for frequently accessed text

### 3. Error Handling

- Always handle file loading errors gracefully
- Provide fallback content for missing assets
- Log errors appropriately without crashing the application

### 4. State Management

- Keep game state separate from rendering logic
- Use immutable updates where possible
- Handle terminal resize events properly

### 5. Input Processing

- Provide multiple key bindings for common actions (WASD + arrows)
- Handle edge cases like rapid key presses
- Implement consistent quit mechanisms

### 6. Animation Guidelines

- Keep animation files organized in the assets directory
- Use descriptive frame separators
- Test animations at different frame rates
- Provide static fallbacks for animation failures

### 7. Internationalization

- Plan your key naming convention early
- Keep placeholder order consistent
- Test with different string lengths
- Provide meaningful fallback keys

This documentation provides a comprehensive guide to using the Project Red Engine for creating terminal-based games. The engine's architecture supports complex game logic while maintaining simplicity in the core APIs.