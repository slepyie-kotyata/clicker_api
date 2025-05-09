package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"math"

	"gorm.io/gorm"
)

func (s *Session) PassiveSellUpdate(upgrade_stats service.UpgradeStats, seconds uint, current_prestige float64) {
	if s.Session.Dishes <= 0 || s.Session.Dishes < 3 {
		return 
	}

	var (
		total_money_per_second float64 = upgrade_stats.MoneyPerSecond
		total_money_passive_multiplier float64 = upgrade_stats.PassiveMoneyMultiplier
		total_sold_per_sell float64 = upgrade_stats.SoldPerSell
	)

	if total_money_per_second == 0 && total_money_passive_multiplier == 0 && total_sold_per_sell == 0{
		return
	} else {

		if total_money_per_second == 0 {
			total_money_per_second = 1
		}

		if total_money_passive_multiplier == 0 {
			total_money_passive_multiplier = 1
		}

		if total_sold_per_sell == 0 {
			total_sold_per_sell = 1
		}
	}

	database.DB.Model(&s.Session).Updates(map[string]interface{}{
		"money": gorm.Expr("money + ?", uint(math.Ceil(total_money_per_second * total_money_passive_multiplier * float64(seconds) * current_prestige * total_sold_per_sell))),
		"dishes": gorm.Expr("dishes - ?", 1 * seconds),
	})
	database.DB.Model(&models.Level{}).Where("session_id = ?", s.Session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 0.05 * float64(seconds) * total_money_per_second))
}

func (s *Session) PassiveCookUpdate(upgrade_stats service.UpgradeStats, seconds uint, current_prestige float64) {
	var (
		total_dishes_per_second float64 = upgrade_stats.DishesPerSecond
		total_dishes_passive_multiplier float64 = upgrade_stats.PassiveDishesMultiplier
	)

	if total_dishes_passive_multiplier == 0 && total_dishes_per_second == 0{
		return
	} else {

		if total_dishes_per_second == 0 {
			total_dishes_per_second = 1
		}

		if total_dishes_passive_multiplier == 0 {
			total_dishes_passive_multiplier = 1
		}
	}

	database.DB.Model(&s.Session).Update("dishes", gorm.Expr("dishes + ?", uint(math.Ceil(total_dishes_per_second * total_dishes_passive_multiplier * float64(seconds) * current_prestige))))
	database.DB.Model(&models.Level{}).Where("session_id = ?", s.Session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 0.2 * float64(seconds) * total_dishes_per_second))
}

func (s *Session) PrestigeUpgrade (upgrade_stats service.UpgradeStats, seconds uint) {
	var (
		total_money_per_second float64 = upgrade_stats.MoneyPerSecond
		total_money_passive_multiplier float64 = upgrade_stats.PassiveMoneyMultiplier
	)

	if total_money_per_second == 0 && total_money_passive_multiplier == 0 {
		return
	} else {

		if total_money_per_second == 0 {
			total_money_per_second = 1
		}

		if total_money_passive_multiplier == 0 {
			total_money_passive_multiplier = 1
		}
	}

	d := total_money_per_second * total_money_passive_multiplier
	p := (d / 10000) * float64(seconds)
	p = math.Round(p * 10000) / 10000

	database.DB.Model(&models.Prestige{}).Where("session_id = ?", s.Session.ID).Update("current_value", gorm.Expr("current_value + ?", p))
}