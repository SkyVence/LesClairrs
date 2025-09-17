package engine

import "sync"

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
