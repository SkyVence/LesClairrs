# UI actions, input handling, and updating the view

This document contains concrete examples for common tasks: handling user input, performing an in-game action (e.g., attack or move), updating animations, and returning `ui.Cmd` side-effects.

Key concepts
-
- `engine.Msg` — messages representing events (key presses, window resize, animation tick, network responses).
- `engine.Cmd` — a function that performs side-effects (timers, IO) and may send messages back to the program.
- `Model.Update(msg)` — receives messages, updates state, and returns an updated model and a command.
- `Model.View()` — returns the view string to render the current state.

Example 1 — simple key input to change screens
-
When the player presses Enter on the menu, the `game` should switch to the game state.

Pseudocode for `Update` handling the menu selection:

```go
func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg := msg.(type) {
    case engine.KeyMsg:
        if msg.Rune == '\r' || msg.Rune == '\n' {
            chosen := m.menu.GetSelected()
            switch chosen.Value {
            case "start":
                m.state = stateGame
                return m, nil
            case "quit":
                return m, engine.Quit
            }
        }
    }
    return m, nil
}
```

Example 2 — triggering an in-game action (attack) with animation
-
This example shows how to handle an attack action, launch an animation command, and transition state when the animation finishes.

Assumptions:
- `engine.KeyMsg` carries the key rune for key presses.
- `engine.Cmd` is a function that can send a message back to the `Update` loop when finished.
- `player.Run()` returns a `engine.Cmd` that sends `RunDoneMsg` when the animation completes.

```go
// Declare a message that indicates attack finished
type RunDoneMsg struct{ Success bool }

func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg := msg.(type) {
    case engine.KeyMsg:
        if msg.Rune == 'r' && m.state == stateGame {
            // Update state immediately to playing-run animation
            m.player.StartRunAnimation()
            // Return a Cmd that will send RunDoneMsg when the animation completes
            return m, m.player.Run()
        }

    case RunDoneMsg:
        if msg.Success {
            // apply side effects of running, etc.
        }
        // return nil - no further commands needed right now
        return m, nil
    }

    return m, nil
}
```

Example 3 — animations driven by a ticker command
-
Animations often need periodic ticks. The `engine` package provides an `engine.Tick` command that starts a goroutine sending `TickMsg{}` every N milliseconds.

```go
func (m *model) Update(msg engine.Msg) (engine.Model, engine.Cmd) {
    switch msg.(type) {
    case engine.TickMsg:
        m.player.AdvanceFrame()
        return m, engine.Tick(100 * time.Millisecond)
    }
    return m, nil
}
```

Example 4 — combining commands
-
You may want to run multiple commands in response to a single message. The `engine` package does not provide a helper for this, but you can implement one like this:

```go
func Batch(cmds ...engine.Cmd) engine.Cmd {
    return func() engine.Msg {
        for _, cmd := range cmds {
            if cmd != nil {
                go func() {
                    // We assume that the commands will send messages to the program loop.
                    // This implementation does not collect the messages.
                    cmd()
                }()
            }
        }
        return nil
    }
}
```

Edge cases and tips
-
- Input during animations: be explicit about whether input is accepted during certain animations (ignore or buffer input).
- Long-running commands: ensure they don't leak goroutines; provide cancellation messages or contexts when appropriate.
- Resize events: update layout/positions in response to `SizeMsg` so UI always fits terminal size.
- Testing: prefer returning deterministic commands in tests (or use mock `send` functions) and assert state transitions and `View()` output.

See also
-
- `docs/game-loop.md` for the overall architecture and message flow.
- `docs/tutoriel-ui.md` for higher-level tutorials and examples.
