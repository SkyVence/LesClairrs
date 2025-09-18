# 🔧 Engine Developer API Documentation

A comprehensive guide for developers who want to understand or extend the ProjectRed terminal game engine.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Core Architecture](#core-architecture)
3. [Program Lifecycle](#program-lifecycle)
4. [Rendering System](#rendering-system)
5. [Input Handling](#input-handling)
6. [Animation System](#animation-system)
7. [Internationalization](#internationalization)
8. [Message System](#message-system)
9. [Engine Configuration](#engine-configuration)
10. [Extending the Engine](#extending-the-engine)
11. [Best Practices](#best-practices)

---

## Overview

The ProjectRed Engine is a terminal-based game engine built in Go, designed around the Model-View-Update (MVU) architecture. This documentation covers the low-level engine components that handle terminal control, rendering, input processing, and application lifecycle.

### Engine Components

- **Program Management**: Application lifecycle and terminal setup
- **Rendering System**: Terminal output control and screen management
- **Input System**: Keyboard input processing and message generation
- **Animation Framework**: Frame-based animation system
- **Internationalization**: Multi-language support with JSON catalogs
- **Message System**: Event-driven communication between components

---

## Core Architecture

### Game Interface

The engine defines a minimal interface that all games must implement:

```go
package engine

type Game interface {
    Init() Msg                    // Initialize game state
    Update(Msg) (Model, Cmd)     // Update state based on messages
    View() string                // Render current state as string
}
```

### Engine Model Wrapper

The engine wraps games to make them compatible with the internal Model interface:

```go
type engineModel struct {
    game Game
}

func Wrap(g Game) Model {
    return &engineModel{game: g}
}
```

**Usage Example:**
```go
game := mypackage.NewMyGame()
wrappedModel := engine.Wrap(game)
program := engine.NewProgram(wrappedModel)
```

---

## Program Lifecycle

### Program Creation and Configuration

```go
func NewProgram(model Model, opts ...ProgramOption) *Program
```

**Available Options:**
- `WithAltScreen()`: Enable alternate screen buffer for full-screen applications

**Example:**
```go
program := engine.NewProgram(
    engine.Wrap(game),
    engine.WithAltScreen(),
)
```

### Terminal Size Detection

```go
func (p *Program) GetSize() (width, height int)
```

The engine automatically detects terminal dimensions and sends `SizeMsg` when the terminal is resized.

### Running the Program

```go
func (p *Program) Run() error
```

This method:
1. Sets up terminal state (raw mode, alternate screen)
2. Starts the rendering system
3. Begins input processing
4. Runs the main event loop
5. Cleans up on exit

---

## Rendering System

### Renderer Interface

```go
type Renderer interface {
    Start()                          // Initialize renderer
    Stop()                          // Graceful shutdown
    Kill()                          // Force shutdown
    Write(string)                   // Output content
    ClearScreen()                   // Clear display
    Repaint()                       // Force full redraw
    ShowCursor()                    // Make cursor visible
    HideCursor()                    // Hide cursor
    SetWindowTitle(string)          // Set terminal title
    AltScreen() bool               // Check alt screen status
    EnterAltScreen()               // Enable alt screen
    ExitAltScreen()                // Disable alt screen
    SetCursor(x, y int)            // Position cursor
    GetSize() (width, height int)  // Get terminal dimensions
}
```

### Standard Renderer

The `StandardRenderer` implementation provides:

- **Buffered Output**: Reduces flickering by batching writes
- **Frame Rate Control**: Configurable rendering speed
- **Terminal Control**: ANSI escape sequence management
- **Cursor Management**: Position and visibility control

**Key Methods:**
```go
renderer := &StandardRenderer{
    frameRate: 60 * time.Millisecond,
}

renderer.Write("Hello, World!")
renderer.SetCursor(10, 5)
renderer.HideCursor()
```

---

## Input Handling

### Input Processing

The engine reads raw terminal input and converts it to structured messages:

```go
func ReadInput(msgs chan<- Msg)
```

### Supported Input Types

**Arrow Keys:**
- `↑` (Up Arrow): `KeyMsg{Rune: '↑'}`
- `↓` (Down Arrow): `KeyMsg{Rune: '↓'}`
- `←` (Left Arrow): `KeyMsg{Rune: '←'}`
- `→` (Right Arrow): `KeyMsg{Rune: '→'}`

**Special Keys:**
- `Ctrl+C`: `QuitMsg{}`
- Regular characters: `KeyMsg{Rune: rune}`

**Implementation Details:**
```go
// Handle escape sequences (arrow keys)
if len(data) >= 3 && data[0] == 0x1b && data[1] == '[' {
    switch data[2] {
    case 'A': msgs <- KeyMsg{Rune: '↑'}
    case 'B': msgs <- KeyMsg{Rune: '↓'}
    case 'C': msgs <- KeyMsg{Rune: '→'}
    case 'D': msgs <- KeyMsg{Rune: '←'}
    }
}
```

---

## Animation System

### Animation Loading

```go
func LoadAnimationFile(filename string) ([]string, error)
```

Loads frame-based animations from text files with `---` separators.

**File Format:**
```
Frame 1 content
Multiple lines supported

---

Frame 2 content
With ASCII art
  ┌─────┐
  │ Art │
  └─────┘

---

Frame 3 content
```

### Animation Structure

```go
type Animation struct {
    Frames []string      // Animation frames
    frame  int          // Current frame index
    speed  time.Duration // Frame duration
}

func NewAnimation(frames []string) Animation
func (a *Animation) Next() string
func (a *Animation) SetSpeed(d time.Duration)
```

**Usage Example:**
```go
frames, err := engine.LoadAnimationFile("assets/animations/loading.txt")
if err != nil {
    log.Fatal(err)
}

animation := engine.NewAnimation(frames)
animation.SetSpeed(200 * time.Millisecond)

// In render loop
currentFrame := animation.Next()
```

---

## Internationalization

### Language Catalog System

```go
type Catalog map[string]string

func Load(lang string) (Catalog, error)
func (c Catalog) Text(key string, args ...any) string
```

### Language File Structure

**File Location:** `assets/interface/{language}.json`

**Example JSON:**
```json
{
  "ui": {
    "title": "ProjectRed RPG",
    "menu": {
      "start": "Start Game",
      "quit": "Quit"
    }
  },
  "game": {
    "health": "Health: {hp}/{max}",
    "level": "Level {level}"
  }
}
```

### Text Interpolation

```go
catalog, err := engine.Load("en")
if err != nil {
    log.Fatal(err)
}

// Simple text
title := catalog.Text("ui.title")

// With placeholders
health := catalog.Text("game.health", currentHP, maxHP)
level := catalog.Text("game.level", playerLevel)
```

**Placeholder System:**
- Placeholders: `{name}`, `{hp}`, `{level}`, etc.
- Arguments replaced in encounter order
- Missing keys return: `⟦key.name⟧`

---

## Message System

### Core Message Types

```go
type Msg interface{}

// User input
type KeyMsg struct {
    Rune rune  // The character or special key
}

// Application control
type QuitMsg struct{}

// Terminal events
type SizeMsg struct {
    Width  int
    Height int
}

// Timer events
type TickMsg struct {
    Time time.Time
}
```

### Command System

```go
type Cmd func() Msg

// Utility commands
func Quit() Msg
func Tick(d time.Duration) Cmd
func TickNow() Cmd
```

**Example Usage:**
```go
func (g *MyGame) Update(msg Msg) (Model, Cmd) {
    switch msg := msg.(type) {
    case KeyMsg:
        if msg.Rune == 'q' {
            return g, engine.Quit
        }
    case TickMsg:
        // Handle timer event
        return g, engine.Tick(1 * time.Second)
    }
    return g, nil
}
```

---

## Engine Configuration

### Engine Configuration Structure

```go
type EngineConfig struct {
    UseAltScreen     bool  // Enable alternate screen buffer
    TargetFPS        int   // Rendering frame rate
    EnableDebugMode  bool  // Debug output enabled
    MaxMessageBuffer int   // Input message buffer size
}

func DefaultEngineConfig() EngineConfig
```

**Default Configuration:**
```go
config := EngineConfig{
    UseAltScreen:     true,
    TargetFPS:        60,
    EnableDebugMode:  false,
    MaxMessageBuffer: 100,
}
```

---

## Extending the Engine

### Creating Custom Renderers

Implement the `Renderer` interface:

```go
type CustomRenderer struct {
    // Your fields
}

func (r *CustomRenderer) Start() {
    // Initialization logic
}

func (r *CustomRenderer) Write(content string) {
    // Custom output logic
}

// Implement all other Renderer methods...
```

### Adding New Message Types

```go
type CustomMsg struct {
    Data string
}

// In your game's Update method
func (g *MyGame) Update(msg Msg) (Model, Cmd) {
    switch msg := msg.(type) {
    case CustomMsg:
        // Handle custom message
        return g, nil
    default:
        // Handle engine messages
        return g.handleEngineMsg(msg)
    }
}
```

### Custom Input Processing

While the engine handles basic input, you can extend it:

```go
func CustomInputProcessor(msgs chan<- engine.Msg) {
    // Your custom input logic
    // Send custom messages to the channel
    msgs <- CustomMsg{Data: "custom event"}
}
```

---

## Best Practices

### 1. Terminal Management

- Always use `WithAltScreen()` for full-screen applications
- Handle terminal resize events (`SizeMsg`) properly
- Clean up terminal state on exit

### 2. Rendering Optimization

- Minimize string allocations in `View()` methods
- Use string builders for complex output
- Avoid unnecessary redraws

```go
func (g *MyGame) View() string {
    var sb strings.Builder
    sb.WriteString("Static content")
    sb.WriteString(fmt.Sprintf("Dynamic: %d", g.value))
    return sb.String()
}
```

### 3. Input Handling

- Provide multiple key bindings for actions (WASD + arrows)
- Handle edge cases like rapid key presses
- Implement consistent quit mechanisms

### 4. Animation Guidelines

- Keep animation files in `assets/animations/`
- Use descriptive frame separators (`---`)
- Test animations at different speeds
- Provide static fallbacks for missing files

### 5. Internationalization

- Plan key naming conventions early
- Keep placeholder order consistent
- Test with different string lengths
- Provide meaningful fallback keys

### 6. Error Handling

- Handle file loading errors gracefully
- Provide fallback content for missing assets
- Log errors without crashing the application

### 7. Performance

- Use buffered rendering for smooth output
- Limit animation frame rates appropriately
- Cache language strings for frequent use

---

This documentation provides comprehensive coverage of the ProjectRed Engine's internal systems, enabling developers to understand, extend, and optimize the engine for their specific needs.