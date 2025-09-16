# Game loop and architecture

This document explains the separation between the engine (program loop) and the game logic in this project, the message/update/view flow, and a minimal code example showing how to adapt a `Game` implementation to the `ui` program.

Goals:
- Keep the engine (program loop, input/event delivery, tick scheduling) independent of game logic.
- Make `game` implement a small `Game` interface so it can be plugged into the engine and tested in isolation.
- Keep UI rendering responsibilities in the `ui` package; the game returns view strings and commands.

Overview
-
The project splits responsibilities across three layers:

1. `ui` package
   - Provides the terminal UI runtime and `Program` that runs an event loop.
   - Defines `Msg` (messages/events), `Cmd` (side-effect functions), and the `Model` interface with `Init() Msg`, `Update(Msg) (Model, Cmd)`, and `View() string`.

2. `engine` package
   - Small adapter that wraps a `Game` implementation and exposes the `ui.Model` interface.
   - Responsible for wiring the game's lifecycle into the `ui.Program` loop but intentionally keeps logic minimal.

3. `game` package
   - The game implementation: state, world, player, enemies, menus, and game rules.
   - Implements the `Game` interface expected by `engine.Wrap(...)`.

Message flow (high level)
-
1. `ui.Program` receives an external event (key press, window resize, timer tick). The event is converted to a `ui.Msg`.
2. `ui.Program` calls the `Model.Update(msg)` function. In our architecture that `Model` is the `engine` adapter which delegates to the `game`.
3. `game.Update(msg)` mutates game state if needed, schedules `ui.Cmd` side-effects (such as playing sounds or starting animations), and returns the updated model and a command.
4. The `ui.Program` executes the returned `ui.Cmd` (asynchronously) and uses the returned model's `View()` to render the screen.

Design notes
-
- The `engine` should not contain game rules. It only adapts the `Game` to `ui.Model`.
- `game` may import `ui` types (messages, commands, components) to build and return view strings and commands.
- Keep `game`'s public surface small; expose a constructor `NewGame()` which returns a `Game` implementation.

Minimal adapter example
-
Here's an illustrative (simplified) adapter used by the engine to adapt a `Game` to the `ui` runtime:

```go
package engine

import "projectred-rpg.com/ui"

// Game is a minimal interface a game package must implement to be run by the engine.
type Game interface {
    Init() ui.Msg
    Update(ui.Msg) (ui.Model, ui.Cmd)
    View() string
}

// engineModel adapts the Game to the ui.Model interface.
type engineModel struct {
    game Game
}

func Wrap(g Game) ui.Model { return &engineModel{game: g} }

func (e *engineModel) Init() ui.Msg    { return e.game.Init() }
func (e *engineModel) Update(m ui.Msg) (ui.Model, ui.Cmd) { return e.game.Update(m) }
func (e *engineModel) View() string   { return e.game.View() }
```

Threading and timing
-
- All state updates happen synchronously inside `Update` and should be deterministic.
- For background/async work create `ui.Cmd` functions that run concurrently; when they complete they should send messages back into the program using the program's message channel (the `ui` runtime provides helpers for this pattern).

Tips
-
- Keep `Update` small and focused: handle the message, change minimal state, return any `Cmd` needed to perform side-effects.
- Return small, composable `Cmd`s. Prefer `Cmd` that send a message to the `Update` loop instead of mutating state externally.
- For animations, return a `Cmd` that drives the animation (e.g., a ticker) and sends frame messages back.
- Write unit tests around `game.Update` and `game.View()` — they are pure/deterministic by design and easy to test.

See also
-
- `docs/ui-actions.md` (examples for specific UI actions and input handling)
- `docs/tutoriel-ui.md` (tutorials and higher-level documentation)
