package types

// TileMap represents a simple ASCII tile map to render in the game area.
// Each string in Tiles represents one row; each rune is rendered as-is.
type TileMap struct {
	Width  int
	Height int
	Tiles  [][]rune
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
	return tm
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
