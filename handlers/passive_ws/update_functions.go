package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"math"

	"gorm.io/gorm"
)

func (s *Session) PassiveSellUpdate(upgrade_stats service.UpgradeStats, seconds uint, prestige_boost float64) {
	if s.Session.Dishes <= 0 || s.Session.Dishes < 3 {
		return 
	}

	if upgrade_stats.MpS == 0 {
		return
	}

	service.SetDefaults(&upgrade_stats)

	minNum := min((float64(seconds) * upgrade_stats.SpS), float64(s.Session.Dishes))

	if s.Session.Level.Rank < 100 {
		database.DB.Model(&models.Level{}).Where("session_id = ?", s.Session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", math.Abs(0.05 * float64(seconds) * upgrade_stats.MpS)))
	}

	database.DB.Model(&s.Session).Updates(map[string]interface{}{
		"money": gorm.Expr("money + ?", uint(math.Ceil(upgrade_stats.MpS * upgrade_stats.MpM * float64(seconds) * prestige_boost * minNum))),
		"dishes": gorm.Expr("dishes - ?", uint(math.Ceil(minNum))),
	})
}

func (s *Session) PassiveCookUpdate(upgrade_stats service.UpgradeStats, seconds uint, prestige_boost float64) {
	if upgrade_stats.DpS == 0 {
		return
	}

	service.SetDefaults(&upgrade_stats)

	if s.Session.Level.Rank < 100 {
		database.DB.Model(&models.Level{}).Where("session_id = ?", s.Session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", math.Abs(0.2 * float64(seconds) * upgrade_stats.DpS)))
	}

	database.DB.Model(&s.Session).Update("dishes", gorm.Expr("dishes + ?", uint(math.Ceil(upgrade_stats.DpS * upgrade_stats.DpM * float64(seconds) * prestige_boost))))
}

func (s *Session) PrestigeUpgrade (upgrade_stats service.UpgradeStats, seconds uint) {
	if upgrade_stats.MpS == 0 {
		return
	}

	service.SetDefaults(&upgrade_stats)

	d := upgrade_stats.MpS * upgrade_stats.MpM
	p := (d / 10000) * float64(seconds)
	p = math.Round(p * 10000) / 10000

	database.DB.Model(&models.Prestige{}).Where("session_id = ?", s.Session.ID).Update("current_value", gorm.Expr("current_value + ?", p))
}