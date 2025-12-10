package ws

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"math"

	"github.com/dariubs/percent"
)

func (s *SessionConn) Buy(id uint) (map[string]interface{}, RequestType) {
	var (
		level_xp     models.LevelXP
		this_upgrade service.FilteredUpgrade
		result_price uint = 0
		exist bool = false
		xp_increase  float64
	)

	upgrade_id := id
	s.session = database.GetSessionState(s.user_id)
	
	if s.session.LevelRank == 100 {
		xp_increase = 0
	} else {
		database.DB.Where("rank = ?", s.session.LevelRank + 1).Find(&level_xp)
		xp_increase = percent.Percent(5, int(level_xp.XP))
	}

	for _, upgrade := range service.FilterUpgrades(s.session, false) {
		if upgrade.ID == upgrade_id {
			this_upgrade = upgrade
			exist = true
		}
	}

	if !exist {
		return map[string]interface{}{
			"message": "upgrade not found",
		}, ErrorRequest
	}

	if this_upgrade.TimesBought == 0 {
		result_price = this_upgrade.Price
	} else {
		result_price = uint(math.Round(float64(this_upgrade.Price) * math.Pow(this_upgrade.PriceFactor, float64(this_upgrade.TimesBought))))
	}

	if s.session.Money < result_price {
		return map[string]interface{}{
			"message": "not enough money, you have: " + utils.IntToString(int(s.session.Money)) + ", you need: " + utils.IntToString(int(result_price)),
		}, ErrorRequest
	}

	s.session.Money -= result_price
	s.session.LevelXP = math.Round((s.session.LevelXP + xp_increase) * 100) / 100
	s.session.Upgrades[upgrade_id] += 1

	database.SaveSessionState(s.user_id, s.session)

	return map[string]interface{}{
		"money": s.session.Money,
		"xp": s.session.LevelXP,
	}, BuyRequest
}

func (s *SessionConn) Cook() (map[string]interface{}, RequestType) {
	s.session = database.GetSessionState(s.user_id)

	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(s.session, true))
	service.SetDefaults(&upgrade_stats)
	
	if !upgrade_stats.HasDish {
		return map[string]interface{}{
			"message": "can't perform this action",
		}, ErrorRequest
	}

	s.session.Dishes += uint(math.Ceil((1 + upgrade_stats.DpC) * upgrade_stats.Dm))
	s.session.LevelXP = math.Round((s.session.LevelXP + 10) * 100) / 100

	database.SaveSessionState(s.user_id, s.session)

	return map[string]interface{}{
		"dishes": s.session.Dishes,
		"xp": s.session.LevelXP,
	}, CookRequest
}

func (s *SessionConn) Sell() (map[string]interface{}, RequestType) {
	s.session = database.GetSessionState(s.user_id)

	if s.session.Dishes <= 0 {
		return map[string]interface{}{
			"message": "not enough dishes",
		}, ErrorRequest
	}

	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(s.session, true))
	service.SetDefaults(&upgrade_stats) 

	min_num := min(upgrade_stats.SpS, float64(s.session.Dishes))
	
	prestige_boost := s.session.PrestigeAccumulated
	if prestige_boost == 0 {
		prestige_boost = 1
	}

	s.session.Money += uint(math.Ceil(upgrade_stats.MpC * upgrade_stats.Mm * min_num * prestige_boost))
	s.session.Dishes -= uint(min_num)
	s.session.LevelXP = math.Round((s.session.LevelXP + 10) * 100) / 100

	database.SaveSessionState(s.user_id, s.session)

	return map[string]interface{}{
		"dishes": s.session.Dishes,
		"money": s.session.Money,
		"xp": s.session.LevelXP,
	}, SellRequest
}

func (s *SessionConn) ListUpgrades() (map[string]interface{}, RequestType) {
	s.session = database.GetSessionState(s.user_id)

	return map[string]interface{}{
		"upgrades": service.FilterUpgrades(s.session, false),
	}, ListRequest
}

func (s *SessionConn) LevelUp() (map[string]interface{}, RequestType) {
	var next_level models.LevelXP

	s.session = database.GetSessionState(s.user_id)

	if s.session.LevelRank == 100 {
		return map[string]interface{}{
			"current_rank": s.session.LevelRank,
			"current_xp": s.session.LevelXP,
		}, LevelUpRequest
	}

	database.DB.Where("rank = ?", s.session.LevelRank + 1).First(&next_level)

	if s.session.LevelXP == float64(next_level.XP) {
		s.session.LevelRank += 1
		s.session.LevelXP = 0

		var new_next_level models.LevelXP
		database.DB.Where("rank = ?", s.session.LevelRank + 1).First(&new_next_level)

		database.SaveSessionState(s.user_id, s.session)
		return map[string]interface{}{
			"current_rank": s.session.LevelRank,
			"current_xp":   s.session.LevelXP,
			"next_xp":      new_next_level.XP,
		}, LevelUpRequest
	}

	if s.session.LevelXP > float64(next_level.XP) {
		s.session.LevelXP = math.Round((s.session.LevelXP - float64(next_level.XP)) * 100) / 100
		s.session.LevelRank += 1

		var new_next_level models.LevelXP
		database.DB.Where("rank = ?", s.session.LevelRank + 1).First(&new_next_level)

		database.SaveSessionState(s.user_id, s.session)
		return map[string]interface{}{
			"current_rank": s.session.LevelRank,
			"current_xp":   s.session.LevelXP,
			"next_xp":      new_next_level.XP,
		}, LevelUpRequest
	}

	database.SaveSessionState(s.user_id, s.session)
	return map[string]interface{}{
		"current_rank": s.session.LevelRank,
		"current_xp":   s.session.LevelXP,
	}, LevelUpRequest
}

func (s *SessionConn) GetLevel() (map[string]interface{}, RequestType) {
	var level models.LevelXP

	s.session = database.GetSessionState(s.user_id)
	if s.session.LevelRank == 100 {
		return map[string]interface{}{
			"current_rank": s.session.LevelRank,
			"current_xp":   s.session.LevelXP,
		}, CheckLevelRequest
	}

	database.DB.Where("rank = ?", s.session.LevelRank + 1).First(&level)

	return map[string]interface{}{
		"current_rank": s.session.LevelRank,
		"current_xp":   s.session.LevelXP,
		"needed_xp":    level.XP,
	}, CheckLevelRequest
}
