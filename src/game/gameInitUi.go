package game

import (
	"projectred-rpg.com/engine"
	"projectred-rpg.com/game/types"
	"projectred-rpg.com/ui"
)

func InitMainMenu(locManager *engine.LocalizationManager) ui.Menu {
	menuOptions := []ui.MenuOption{
		{Label: locManager.Text("ui.menu.start"), Value: "start"},
		{Label: locManager.Text("ui.menu.settings"), Value: "settings"},
		{Label: locManager.Text("ui.menu.quit"), Value: "quit"},
	}
	menu, err := ui.NewMenuWithArtFromFile(locManager.Text("ui.menu.mainmenu"), menuOptions, "assets/logo.txt")
	if err != nil {
		menu = ui.NewMenu("ui.menu.mainmenu", menuOptions)
	}
	return menu
}

func InitializeClassSelection(locManager *engine.LocalizationManager, classes []types.Class) ui.ClassMenu {
	menuOptions := []ui.ClassMenuOption{}
	for _, class := range classes {
		label := class.Name
		if locManager != nil {
			translatedLabel := locManager.Text(class.Name)
			if translatedLabel != "" && translatedLabel != class.Name {
				label = translatedLabel
			}
		}
		menuOptions = append(menuOptions, ui.ClassMenuOption{
			Label: label,
			Value: class.Name,
			Desc:  class.Description,
			Stats: ui.ClassStats{
				MaxHP:    class.MaxHP,
				Force:    class.Force,
				Speed:    class.Speed,
				Defense:  class.Defense,
				Accuracy: class.Accuracy,
			},
		})
	}
	menu := ui.NewClassMenu(locManager.Text("ui.class.menu.name"), menuOptions, locManager)
	return menu
}

func InitializeSettingsSelection(locManager *engine.LocalizationManager, languageOptions []string) ui.SettingsMenu {
	menuOptions := []ui.SettingsMenuOption{}
	for _, lang := range languageOptions {
		menuOptions = append(menuOptions, ui.SettingsMenuOption{
			Label: lang, // Keep the raw language code for localization lookup
			Value: lang,
		})
	}
	menu := ui.NewSettingsMenu(locManager.Text("ui.settings.menu.title"), menuOptions, locManager)
	return menu
}

func InitializeMerchantMenu(locManager *engine.LocalizationManager) ui.MerchantMenu {
	menuOptions := []ui.MerchantMenuOption{
		// === WEAPONS SECTION ===
		{
			Item: types.Item{
				Name:        "═══ WEAPONS ═══",
				Description: "",
				Type:        types.Weapon,
			},
			Price: 0,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.katana.name",
				Description: "ui.weapons.katana.description",
				Type:        types.Weapon,
			},
			Price: 150,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.faux neuralink.name",
				Description: "ui.weapons.faux neuralink.description",
				Type:        types.Weapon,
			},
			Price: 200,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.arc synaptique.name",
				Description: "ui.weapons.arc synaptique.description",
				Type:        types.Weapon,
			},
			Price: 175,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.sniper.name",
				Description: "ui.weapons.sniper.description",
				Type:        types.Weapon,
			},
			Price: 250,
		},
		{
			Item: types.Item{
				Name:        "ui.weapons.neon reaver.name",
				Description: "ui.weapons.neon reaver.description",
				Type:        types.Weapon,
			},
			Price: 300,
		},
		{
			Item: types.Item{
				Name:        "═══ CONSUMABLE ═══",
				Description: "",
				Type:        types.Consumable,
			},
			Price: 0,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.small_medkit.name",
				Description: "ui.consumable.small_medkit.description",
				Type:        types.Consumable,
			},
			Price: 25,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.large_medkit.name",
				Description: "ui.consumable.large_medkit.description",
				Type:        types.Consumable,
			},
			Price: 50,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.money.name",
				Description: "ui.consumable.money.description",
				Type:        types.Consumable,
			},
			Price: 75,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.serum.name",
				Description: "ui.consumable.serum.description",
				Type:        types.Consumable,
			},
			Price: 100,
		},
		{
			Item: types.Item{
				Name:        "ui.consumable.flash.name",
				Description: "ui.consumable.flash.description",
				Type:        types.Consumable,
			},
			Price: 80,
		},
	}

	merchantName := locManager.Text("game.merchants.weapon")
	menu := ui.NewMerchantMenu(merchantName, menuOptions, locManager)
	return menu
}
