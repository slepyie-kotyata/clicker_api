package service

import (
	"clicker_api/pkg/models"
	"clicker_api/services/main_api/database"
	"sort"
)

type UpgradeStats struct {
	MpS		float64
	MpM 	float64
	SpS     float64
	DpS     float64
	DpC     float64
	Dm      float64
	Mm 		float64
	MpC		float64
	DpM		float64
	HasDish	bool
}

func SetDefaults(stats *UpgradeStats) {
	defaults := map[string]*float64{
		"dPs": &stats.DpS,
		"mPs": &stats.MpS,
		"dpM": &stats.DpM,
		"mpM": &stats.MpM,
		"mM": &stats.Mm,
		"dM": &stats.Dm,
	}

	for _, ptr := range defaults {
		if *ptr == 0 {
			*ptr = 1
		}
	}
}

func CountBoostValues(filtered_upgrades []FilteredUpgrade) UpgradeStats {
	upgrade_stats := UpgradeStats{}

	for _, upgrade := range filtered_upgrades {
		switch upgrade.Boost.BoostType {
		case "mPs":
			upgrade_stats.MpS += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dPs":
			upgrade_stats.DpS += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "mpM":
			upgrade_stats.MpM += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dpM":
			upgrade_stats.DpM += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "sPs":
			upgrade_stats.SpS += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dPc":
			upgrade_stats.DpC += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "mPc":
			upgrade_stats.MpC += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "dM":
			upgrade_stats.Dm += upgrade.Boost.Value * float64(upgrade.TimesBought)
		case "mM":
			upgrade_stats.Mm += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}

		if upgrade.UpgradeType == "dish" {
			upgrade_stats.HasDish = true
		}
	}	

	upgrade_stats.SpS += 1

	return upgrade_stats
}

var upgrade_priority = map[string]int{
	"dish":      	1,
	"equipment": 	2,
	"global":    	3,
	"staff": 		4,
	"point":     	5,
}

func getUpgradeTypePriority(upgradeType models.UpgradeType) int {
	p := upgrade_priority[string(upgradeType)]

	return p
}

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

func FilterUpgrades(session *models.SessionState, is_bought bool) []FilteredUpgrade {
	filtered_upgrades := make([]FilteredUpgrade, 0)

	for _, upgrade := range *database.Upgrades {
		times_bought, ok := session.Upgrades[upgrade.ID]

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
		p_i := getUpgradeTypePriority(filtered_upgrades[i].UpgradeType)
		p_j := getUpgradeTypePriority(filtered_upgrades[j].UpgradeType)

		if p_i != p_j {
			return p_i < p_j
		}
		
		return filtered_upgrades[i].Price < filtered_upgrades[j].Price
	})

	return filtered_upgrades
}