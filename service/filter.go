package service

import (
	"clicker_api/database"
	"clicker_api/models"
	"sort"
)

type FilteredUpgrade struct {
	ID          uint               `json:"id" gorm:"primary_key"`
	Name        string             `json:"name"`
	IconName    string             `json:"icon_name"`
	UpgradeType models.UpgradeType `json:"upgrade_type"`
	PriceFactor float64            `json:"price_factor"`
	Price       uint               `json:"price"`
	AccessLevel uint               `json:"access_level"`
	Boost       models.Boost       `json:"boost"`
	TimesBought uint               `json:"times_bought"`
}

type UpgradeStats struct {
	MoneyPerSecond 			float64
	PassiveMoneyMultiplier 	float64
	DishesPerSecond         float64
	PassiveDishesMultiplier	float64
	HasDish					bool
}

func FilterUpgrades(session models.Session, is_bought bool) []FilteredUpgrade {
	filtered_upgrades := make([]FilteredUpgrade, 0)

	var session_upgrades []models.SessionUpgrade
	database.DB.Where("session_id = ?", session.ID).Find(&session_upgrades)

	times_bought_map := make(map[uint]uint)
	for _, su := range session_upgrades {
		times_bought_map[su.UpgradeID] = su.TimesBought
	}

	for _, upgrade := range session.Upgrades {
		times_bought, ok := times_bought_map[upgrade.ID]

		this_upgrade := FilteredUpgrade{
			ID:          upgrade.ID,
			Name:        upgrade.Name,
			IconName:    upgrade.IconName,
			UpgradeType: upgrade.UpgradeType,
			PriceFactor: upgrade.PriceFactor,
			Price:       upgrade.Price,
			AccessLevel: upgrade.AccessLevel,
			Boost:       upgrade.Boost,
			TimesBought: times_bought,
		}

		if is_bought {
			if ok && times_bought > 0 {
				filtered_upgrades = append(filtered_upgrades, this_upgrade)
			}
		} else {
			if ok && (times_bought == 0 || upgrade.UpgradeType != "dish") {
				filtered_upgrades = append(filtered_upgrades, this_upgrade)
			}
		}
	}
	sort.Slice(filtered_upgrades, func(i, j int) bool {
		return filtered_upgrades[i].ID < filtered_upgrades[j].ID
	})

	return filtered_upgrades
}

func CountBoostValues(filtered_upgrades []FilteredUpgrade) UpgradeStats {
	upgrade_stats := UpgradeStats{}

	for _, upgrade := range filtered_upgrades {
		switch upgrade.Boost.BoostType {
		case "mPs":
			upgrade_stats.MoneyPerSecond += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dPs":
			upgrade_stats.DishesPerSecond += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "mpM":
			upgrade_stats.PassiveMoneyMultiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dpM":
			upgrade_stats.PassiveDishesMultiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}

		if upgrade.UpgradeType == "dish" {
			upgrade_stats.HasDish = true
		}
	}

	return upgrade_stats
}