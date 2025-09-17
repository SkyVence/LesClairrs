package ui

import "projectred-rpg.com/engine"

type Localizable interface {
	RefreshText()

	SetLanguage(lang string) error
}

type LocalizableComponent struct {
	locManager *engine.LocalizationManager
}

func NewLocalizableComponent() LocalizableComponent {
	return LocalizableComponent{
		locManager: engine.GetLocalizationManager(),
	}
}

func (lc *LocalizableComponent) Translate(key string, args ...any) string {
	return lc.locManager.Text(key, args...)
}
