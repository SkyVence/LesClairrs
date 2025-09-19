package engine

import (
	"os"
	"strings"
	"sync"
)

type LocalizationManager struct {
	currentLang     string
	catalog         Catalog
	fallbackCatalog Catalog
	mutex           sync.RWMutex
}

var (
	globalLocManager *LocalizationManager
	once             sync.Once
)

// GetLocalizationManager returns singleton instance of LocalizationManager
// Initializes with default language "fr" on first call
func GetLocalizationManager() *LocalizationManager {
	once.Do(func() {
		globalLocManager = &LocalizationManager{
			currentLang: "fr",
		}
		globalLocManager.SetLanguage("fr")
	})
	return globalLocManager
}

// SetLanguage loads and sets the catalog for the specified language
// Returns error if language file cannot be loaded
func (lm *LocalizationManager) SetLanguage(lang string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	catalog, err := Load(lang)
	if err != nil {
		return err
	}

	lm.currentLang = lang
	lm.catalog = catalog
	return nil
}

// Text retrieves localized text for key with placeholder replacement
// Returns placeholder notation if key not found or catalog not loaded
func (lm *LocalizationManager) Text(key string, args ...any) string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	if lm.catalog != nil {
		return lm.catalog.Text(key, args...)
	}

	return "⟦" + key + "⟧"
}

// GetCurrentLanguage returns the currently set language code
func (lm *LocalizationManager) GetCurrentLanguage() string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.currentLang
}

// GetSupportedLanguages scans assets/interface directory for available .json language files
// Returns slice of language codes and any directory read error
func (lm *LocalizationManager) GetSupportedLanguages() ([]string, error) {
	interfaceDir := "assets/interface"
	files, err := os.ReadDir(interfaceDir)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			lang := strings.TrimSuffix(file.Name(), ".json")
			languages = append(languages, lang)
		}
	}

	return languages, nil
}
