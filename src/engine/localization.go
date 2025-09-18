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

func GetLocalizationManager() *LocalizationManager {
	once.Do(func() {
		globalLocManager = &LocalizationManager{
			currentLang: "fr",
		}
		// Initialize with default language
		globalLocManager.SetLanguage("fr")
	})
	return globalLocManager
}

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

func (lm *LocalizationManager) Text(key string, args ...any) string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	if lm.catalog != nil {
		return lm.catalog.Text(key, args...)
	}

	// Fallback to key name if no catalog loaded
	return "⟦" + key + "⟧"
}

func (lm *LocalizationManager) GetCurrentLanguage() string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.currentLang
}

func (lm *LocalizationManager) GetSupportedLanguages() ([]string, error) {
	interfaceDir := "assets/interface"
	files, err := os.ReadDir(interfaceDir)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			// Extract language code by removing .json extension
			lang := strings.TrimSuffix(file.Name(), ".json")
			languages = append(languages, lang)
		}
	}

	return languages, nil
}
