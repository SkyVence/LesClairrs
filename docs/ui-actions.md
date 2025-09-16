# UI actions, input handling, and updating the view

This document contains concrete examples for common tasks: handling user input, performing an in-game action (e.g., attack or move), updating animations, and returning `ui.Cmd` side-effects.

Key concepts
-
- `ui.Msg` — messages representing events (key presses, window resize, animation tick, network responses).
- `ui.Cmd` — a function that performs side-effects (timers, IO) and may send messages back to the program.
- `Model.Update(msg)` — receives messages, updates state, and returns an updated model and a command.
- `Model.View()` — returns the view string to render the current state.

Example 1 — simple key input to change screens
-
When the player presses Enter on the menu, the `game` should switch to the game state.

Pseudocode for `Update` handling the menu selection:

```go
func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg := msg.(type) {
    case ui.KeyMsg:
        if msg.Rune == '\r' || msg.Rune == '\n' {
            chosen := m.menu.GetSelected()
            switch chosen.Value {
            case "start":
                m.state = stateGame
                // Return a command that initializes the player's animation cycle
                return m, m.player.Init()
            case "quit":
                return m, ui.Quit
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
- `ui.KeyMsg` carries the key rune for key presses.
- `ui.Cmd` is a function that can send a message back to the `Update` loop when finished.
- `player.Attack()` returns a `ui.Cmd` that sends `AttackDoneMsg` when the animation completes.

```go
// Declare a message that indicates attack finished
type AttackDoneMsg struct{ Success bool }

func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg := msg.(type) {
    case ui.KeyMsg:
        if msg.Rune == 'a' && m.state == stateGame {
            // Update state immediately to playing-attack animation
            m.player.StartAttackAnimation()
            // Return a Cmd that will send AttackDoneMsg when the animation completes
            return m, m.player.Attack()
        }

    case AttackDoneMsg:
        if msg.Success {
            // apply damage, update enemy HP, etc.
            m.applyDamageToTarget(10)
        }
        // return nil - no further commands needed right now
        return m, nil
    }

    return m, nil
}
```

Example 3 — animations driven by a ticker command
-
Animations often need periodic ticks. Implement a `Cmd` that starts a goroutine sending `TickMsg{}` every N milliseconds. Use `Cmd` helpers provided by `ui` when available.

```go
// TickMsg indicates a single animation tick
type TickMsg struct{}

func TickCmd(interval time.Duration) ui.Cmd {
    return func(send func(ui.Msg)) {
        ticker := time.NewTicker(interval)
        go func() {
            for range ticker.C {
                send(TickMsg{})
            }
        }()
    }
}

func (m *model) Update(msg ui.Msg) (ui.Model, ui.Cmd) {
    switch msg.(type) {
    case TickMsg:
        m.player.AdvanceFrame()
        return m, nil
    }
    return m, nil
}
```

Example 4 — combining commands
-
You may want to run multiple commands in response to a single message. Compose them with a helper:

```go
func Batch(cmds ...ui.Cmd) ui.Cmd {
    return func(send func(ui.Msg)) {
        for _, cmd := range cmds {
            if cmd == nil { continue }
            go cmd(send)
        }
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
