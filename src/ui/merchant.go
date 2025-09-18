package ui

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/lipgloss"
    "projectred-rpg.com/engine"
    "projectred-rpg.com/game/types"
)

type MerchantMenuOption struct {
    Item  types.Item
    Price int
}

type MerchantMenu struct {
    Title    string
    Options  []MerchantMenuOption
    Styles   MerchantMenuStyles
    Loc      *engine.LocalizationManager
    selected int
    width    int
    height   int
}

type MerchantMenuStyles struct {
    Title       lipgloss.Style
    Selected    lipgloss.Style
    Normal      lipgloss.Style
    Description lipgloss.Style
    Stats       lipgloss.Style
    Sidebar     lipgloss.Style
}

func DefaultMerchantMenuStyles() MerchantMenuStyles {
    return MerchantMenuStyles{
        Title: lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color("#FAFAFA")).
            Background(lipgloss.Color("#7D56F4")).
            Padding(0, 1).
            MarginBottom(1),
        Selected: lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color("#EE6FF8")).
            Background(lipgloss.Color("#654EA3")).
            Padding(0, 1),
        Normal: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FAFAFA")).
            Padding(0, 1),
        Description: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#C0C0C0")).
            Padding(1, 1).
            MarginTop(1),
        Stats: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FAFAFA")).
            Padding(1, 1).
            MarginTop(1),
        Sidebar: lipgloss.NewStyle().
            Background(lipgloss.Color("#1F1F2E")).
            Padding(0, 1),
    }
}

func NewMerchantMenu(title string, options []MerchantMenuOption, loc *engine.LocalizationManager, styles ...MerchantMenuStyles) MerchantMenu {
    menuStyles := DefaultMerchantMenuStyles()
    if len(styles) > 0 {
        menuStyles = styles[0]
    }

    return MerchantMenu{
        Title:   title,
        Options: options,
        Styles:  menuStyles,
        Loc:     loc,
    }
}

func (m MerchantMenu) localize(s string) string {
    if m.Loc == nil || s == "" {
        return s
    }
    tr := m.Loc.Text(s)
    if strings.HasPrefix(tr, "⟦") && strings.HasSuffix(tr, "⟧") {
        return s
    }
    return tr
}

func (m MerchantMenu) Update(msg engine.Msg) (MerchantMenu, engine.Msg) {
    switch msg := msg.(type) {
    case engine.SizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case engine.KeyMsg:
        switch msg.Rune {
        case '↓':
            if m.selected < len(m.Options)-1 {
                m.selected++
            }
        case '↑':
            if m.selected > 0 {
                m.selected--
            }
        }
    }
    return m, nil
}

func (m MerchantMenu) View() string {
    if m.width == 0 || m.height == 0 {
        return ""
    }

    var menuItems []string

    menuItems = append(menuItems, m.Styles.Title.Render(m.localize(m.Title)))
    for i, option := range m.Options {
        var item string
        itemName := m.localize(option.Item.Name)
        if i == m.selected {
            item = m.Styles.Selected.Render("▶ " + itemName)
        } else {
            item = m.Styles.Normal.Render("  " + itemName)
        }
        menuItems = append(menuItems, item)
    }

    leftColumn := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

    const minTotalForSidebar = 44
    canTrySidebar := m.width >= minTotalForSidebar && len(m.Options) > 0

    gapW := 2
    leftW := m.width * 2 / 5
    if leftW < 18 {
        leftW = 18
    }
    if leftW > m.width {
        leftW = m.width
    }

    leftMargin := 0
    if m.width > leftW {
        leftMargin = (m.width - leftW) / 2
    }

    rightW := 0
    if canTrySidebar {
        availableRight := m.width - (leftMargin + leftW) - gapW
        if availableRight >= 20 {
            rightW = availableRight
        }
    }

    targetH := m.height
    if targetH < 1 {
        targetH = 1
    }

    left := lipgloss.Place(leftW, targetH, lipgloss.Left, lipgloss.Center, leftColumn)

    spacer := lipgloss.NewStyle().Width(leftMargin).Height(targetH).Render("")
    gap := lipgloss.NewStyle().Width(gapW).Height(targetH).Render("")

    var content string
    if rightW > 0 {
        rightContent := m.renderSidebar(m.Options[m.selected], rightW)
        right := lipgloss.Place(rightW, targetH, lipgloss.Left, lipgloss.Center, rightContent)
        content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left, gap, right)
    } else {
        content = lipgloss.JoinHorizontal(lipgloss.Top, spacer, left)
    }

    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Left,
        lipgloss.Center,
        content,
    )
}

func (m MerchantMenu) renderSidebar(opt MerchantMenuOption, width int) string {
    if width <= 0 {
        return ""
    }
    
    itemName := m.localize(opt.Item.Name)
    itemType := m.localize(opt.Item.Description)

    nameBlock := m.Styles.Description.
        Width(width).
        Render(itemName)

    descBlock := m.Styles.Description.
        Width(width).
        Render("Type: " + itemType)

    priceBlock := m.Styles.Stats.
        Width(width).
        Render(fmt.Sprintf("Prix: %d", opt.Price))  // Utilise fmt.Sprintf !

    inner := lipgloss.JoinVertical(lipgloss.Left, nameBlock, descBlock, priceBlock)
    return m.Styles.Sidebar.Width(width).Render(inner)
}

func (m MerchantMenu) GetSelected() MerchantMenuOption {
    if m.selected >= 0 && m.selected < len(m.Options) {
        return m.Options[m.selected]
    }
    return MerchantMenuOption{}
}