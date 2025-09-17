package game

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type WeaponType int

const (
	Melee WeaponType = iota
	Ranged
)

type Weapon struct {
	KeyName string
	Type    WeaponType
	Attacks []Attack
}

type Attack struct {
	KeyName  string
	KeyDesc  string
	Damage   int
	Duration int
	CoolDown int
}

var (
	weaponCache   map[string]Weapon
	weaponMutex   sync.RWMutex
	weaponsLoaded bool = false
)

// LoadWeapons loads all weapons from JSON files in assets/data directory
func LoadWeapons() error {
	weaponMutex.Lock()
	defer weaponMutex.Unlock()

	if weaponsLoaded {
		return nil // Already loaded
	}

	weaponCache = make(map[string]Weapon)

	// Get the path to the assets/data directory
	assetsPath := filepath.Join("assets", "data")

	// Read all JSON files in the directory
	err := filepath.WalkDir(assetsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		// Read the JSON file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read weapon file %s: %w", path, err)
		}

		// Parse the JSON into a Weapon struct
		var weapon Weapon
		if err := json.Unmarshal(data, &weapon); err != nil {
			return fmt.Errorf("failed to parse weapon file %s: %w", path, err)
		}

		// Store the weapon in the cache using its KeyName
		weaponCache[weapon.KeyName] = weapon

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load weapons: %w", err)
	}

	weaponsLoaded = true
	return nil
}

// GetWeapon retrieves a weapon by its key name
func GetWeapon(keyName string) (Weapon, bool) {
	weaponMutex.RLock()
	defer weaponMutex.RUnlock()

	weapon, exists := weaponCache[keyName]
	return weapon, exists
}

// GetAllWeapons returns a copy of all loaded weapons
func GetAllWeapons() map[string]Weapon {
	weaponMutex.RLock()
	defer weaponMutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]Weapon)
	for k, v := range weaponCache {
		result[k] = v
	}
	return result
}

// GetWeaponsByType returns all weapons of a specific type
func GetWeaponsByType(weaponType WeaponType) []Weapon {
	weaponMutex.RLock()
	defer weaponMutex.RUnlock()

	var weapons []Weapon
	for _, weapon := range weaponCache {
		if weapon.Type == weaponType {
			weapons = append(weapons, weapon)
		}
	}
	return weapons
}
