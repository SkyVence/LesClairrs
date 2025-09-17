package ui

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
    "projectred-rpg.com/engine"
)

type ClassCard struct {
    Name        string
    Description string
    MaxHP       int
    Force       int
    Speed       int
    Defense     int
    Accuracy    int
}

type ClassSelection struct {
    Title    string
    Classes  []ClassCard
    Selected int
    width    int
    height   int
    Styles   MenuStyles
}

type ClassChosenMsg struct {
    Class ClassCard
}

type ClassSelectionCanceledMsg struct{}

func NewClassSelection(title string, classes []ClassCard) ClassSelection {
    return ClassSelection{
        Title:   title,
        Classes: classes,
        Styles:  DefaultMenuStyles(),
    }
}

func (c ClassSelection) Init() engine.Msg {
    return nil
}

func (c ClassSelection) Update(msg engine.Msg) (ClassSelection, engine.Cmd) {
    switch m := msg.(type) {
    case engine.SizeMsg:
        c.width = m.Width
        c.height = m.Height
    case engine.KeyMsg:
        switch m.Rune {
        case '↓':
            if c.Selected < len(c.Classes)-1 {
                c.Selected++
            }
        case '↑':
            if c.Selected > 0 {
                c.Selected--
            }
        case '\n', '\r':
            if len(c.Classes) > 0 {
                chosen := c.Classes[c.Selected]
                return c, func() engine.Msg { return ClassChosenMsg{Class: chosen} }
            }
        case 'q':
            return c, func() engine.Msg { return ClassSelectionCanceledMsg{} }
        }
    }
    return c, nil
}

func (c ClassSelection) View() string {
    if c.width == 0 || c.height == 0 {
        return ""
    }

    var lines []string
    lines = append(lines, c.Styles.Title.Render(c.Title))
    lines = append(lines, "")

    // Classes
    for i, class := range c.Classes {
        if i == c.Selected {
            lines = append(lines, c.Styles.Selected.Render("▶ "+class.Name))
        } else {
            lines = append(lines, c.Styles.Normal.Render("  "+class.Name))
        }
    }

    // Stats compactes de la classe sélectionnée
    if len(c.Classes) > 0 {
        sel := c.Classes[c.Selected]
        lines = append(lines, "")
        lines = append(lines, c.Styles.Normal.Render(fmt.Sprintf("HP:%d Force:%d Speed:%d Def:%d Acc:%d", 
            sel.MaxHP, sel.Force, sel.Speed, sel.Defense, sel.Accuracy)))
        lines = append(lines, "")
        lines = append(lines, c.Styles.Normal.Render("[Enter] Select | [q] Back"))
    }

    menu := lipgloss.JoinVertical(lipgloss.Left, lines...)
    return lipgloss.Place(c.width, c.height, lipgloss.Center, lipgloss.Center, menu)
}