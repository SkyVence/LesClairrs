package ui

import (
	"sync"

	"projectred-rpg.com/engine"
)

type ComponentRegistry struct {
	components []Localizable
	mutex      sync.RWMutex
}

var (
	globalRegistry *ComponentRegistry
	registryOnce   sync.Once
)

func GetComponentRegistry() *ComponentRegistry {
	registryOnce.Do(func() {
		globalRegistry = &ComponentRegistry{
			components: make([]Localizable, 0),
		}
	})
	return globalRegistry
}

// Register adds a component to the registry
func (cr *ComponentRegistry) Register(component Localizable) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.components = append(cr.components, component)
}

// RefreshAllComponents updates text for all registered components
func (cr *ComponentRegistry) RefreshAllComponents() {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	for _, component := range cr.components {
		component.RefreshText()
	}
}

// ChangeLanguage changes language for all components
func (cr *ComponentRegistry) ChangeLanguage(lang string) error {
	locManager := engine.GetLocalizationManager()
	if err := locManager.SetLanguage(lang); err != nil {
		return err
	}

	cr.RefreshAllComponents()
	return nil
}
