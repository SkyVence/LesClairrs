package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const levelsDir = "assets/levels"
const levelsDirAlt = "src/assets/levels"

// LoadWorld loads a single world definition from assets/levels/world-[id].json
func LoadWorld(id int) (World, error) {
	var w World
	filename := fmt.Sprintf("world-%d.json", id)
	path := filepath.Join(levelsDir, filename)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			alt := filepath.Join(levelsDirAlt, filename)
			f, err = os.Open(alt)
			if err != nil {
				return w, err
			}
		} else {
			return w, err
		}
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(&w); err != nil {
		return w, err
	}
	for i := range w.Stages {
		w.Stages[i].WorldID = w.WorldID
	}
	return w, nil
}

// LoadLevels loads worlds sequentially starting at 1 until a file is missing.
// It returns all successfully loaded worlds.
func LoadLevels() []World {
	var worlds []World
	for id := 1; ; id++ {
		w, err := LoadWorld(id)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			// Stop on other IO/parse errors as well.
			break
		}
		worlds = append(worlds, w)
	}
	return worlds
}
