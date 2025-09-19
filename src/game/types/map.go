package types

// TransitionZone represents an area that allows transitioning to the next stage/world
type TransitionZone struct {
	X      int
	Y      int
	Width  int
	Height int
	Active bool
}

// IsInZone checks if a position is within the transition zone
func (tz *TransitionZone) IsInZone(x, y int) bool {
	if !tz.Active {
		return false
	}
	return x >= tz.X && x < tz.X+tz.Width && y >= tz.Y && y < tz.Y+tz.Height
}

// TileMap represents a simple ASCII tile map to render in the game area.
// Each string in Tiles represents one row; each rune is rendered as-is.
type TileMap struct {
	Width          int
	Height         int
	Tiles          [][]rune
	TransitionZone *TransitionZone
}

// NewTileMap constructs a TileMap from lines of text.
func NewTileMap(lines []string) *TileMap {
	tm := &TileMap{Height: len(lines)}
	maxW := 0
	tm.Tiles = make([][]rune, len(lines))
	for i, line := range lines {
		runes := []rune(line)
		tm.Tiles[i] = runes
		if len(runes) > maxW {
			maxW = len(runes)
		}
	}
	tm.Width = maxW

	// Create default transition zone in bottom-right corner
	tm.TransitionZone = &TransitionZone{
		X:      maxW - 3,
		Y:      len(lines) - 3,
		Width:  2,
		Height: 2,
		Active: false, // Initially inactive
	}

	return tm
}

// SetCustomTransitionZone sets a custom transition zone with specified coordinates
func (tm *TileMap) SetCustomTransitionZone(x, y, width, height int) {
	if tm == nil {
		return
	}
	tm.TransitionZone = &TransitionZone{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Active: false, // Initially inactive
	}
}

// At returns the rune at x,y if within bounds; otherwise space.
func (tm *TileMap) At(x, y int) rune {
	if tm == nil || y < 0 || y >= len(tm.Tiles) {
		return ' '
	}
	row := tm.Tiles[y]
	if x < 0 || x >= len(row) {
		return ' '
	}
	return row[x]
}

// ActivateTransitionZone enables the transition zone for stage progression
func (tm *TileMap) ActivateTransitionZone() {
	if tm.TransitionZone != nil {
		tm.TransitionZone.Active = true
	}
}

// IsInTransitionZone checks if a player position is in the active transition zone
func (tm *TileMap) IsInTransitionZone(x, y int) bool {
	if tm.TransitionZone == nil {
		return false
	}
	return tm.TransitionZone.IsInZone(x, y)
}
