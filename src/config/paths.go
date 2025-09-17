package config

import "path/filepath"

// AssetPaths contains all asset directory paths
type AssetPaths struct {
	Root          string
	DataDir       string
	AnimationsDir string
	InterfaceDir  string
	LevelsDir     string
	WorldsDir     string
	WeaponsDir    string
	EnemiesDir    string
	ClassesDir    string
}

// DefaultAssetPaths returns the default asset path configuration
func DefaultAssetPaths() AssetPaths {
	root := "assets"
	return AssetPaths{
		Root:          root,
		DataDir:       filepath.Join(root, "data"),
		AnimationsDir: filepath.Join(root, "animations"),
		InterfaceDir:  filepath.Join(root, "interface"),
		LevelsDir:     filepath.Join(root, "levels"),
		WorldsDir:     filepath.Join(root, "levels"),
		WeaponsDir:    filepath.Join(root, "data"),
		EnemiesDir:    filepath.Join(root, "data"),
		ClassesDir:    filepath.Join(root, "data"),
	}
}

// Global asset paths instance
var AssetPathsConfig = DefaultAssetPaths()
