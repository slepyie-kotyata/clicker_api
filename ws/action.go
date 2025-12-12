package ws

import (
	"clicker_api/database"
	"clicker_api/service"
	"clicker_api/utils"
	"log"
	"math"

	"github.com/dariubs/percent"
)

func (s *SessionConn) Buy(id uint) (map[string]interface{}, RequestType) {
	var (
		this_upgrade service.FilteredUpgrade
		result_price uint = 0
		xp_increase  float64
		exist bool = false
	)

	upgrade_id := id
	session := database.GetSessionState(s.user_id)

	log.Println(upgrade_id)
	for _, upgrade := range service.FilterUpgrades(session, false) {
		if upgrade.ID == upgrade_id {
			log.Println()
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

	if session.Money < result_price {
		return map[string]interface{}{
			"message": "not enough money, you have: " + utils.IntToString(int(session.Money)) + ", you need: " + utils.IntToString(int(result_price)),
		}, ErrorRequest
	}

	session.Money -= result_price

	if session.LevelRank == 100 {
		xp_increase = 0
	} else {
		xp_increase = percent.Percent(1, int(database.LevelsXP[session.LevelRank + 1]))
		session.LevelXP = math.Round((session.LevelXP + xp_increase + (session.PrestigeCurrent * 0.05)) * 100) / 100
	}
	
	session.Upgrades[upgrade_id] += 1

	database.SaveSessionState(s.user_id, session)
	database.A.MarkChanged(s.user_id)

	return map[string]interface{}{
		"money": session.Money,
		"xp": session.LevelXP,
	}, BuyRequest
}

func (s *SessionConn) Cook() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)

	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(session, true))
	service.SetDefaults(&upgrade_stats)
	
	if !upgrade_stats.HasDish {
		return map[string]interface{}{
			"message": "can't perform this action",
		}, ErrorRequest
	}

	session.Dishes += uint(math.Ceil((1 + upgrade_stats.DpC) * upgrade_stats.Dm))

	if session.LevelXP != 100 {
		session.LevelXP = math.Round((session.LevelXP + 1 + (session.PrestigeCurrent * 0.05)) * 100) / 100
	}

	database.SaveSessionState(s.user_id, session)
	database.A.MarkChanged(s.user_id)

	return map[string]interface{}{
		"dishes": session.Dishes,
		"xp": session.LevelXP,
	}, CookRequest
}

func (s *SessionConn) Sell() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)

	if session.Dishes <= 0 {
		return map[string]interface{}{
			"message": "not enough dishes",
		}, ErrorRequest
	}

	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(session, true))
	service.SetDefaults(&upgrade_stats) 

	min_num := min(upgrade_stats.SpS, float64(session.Dishes))
	
	prestige_boost := session.PrestigeBoost
	if prestige_boost == 0 {
		prestige_boost = 1
	}

	log.Printf("SpS: %v", upgrade_stats.SpS)
	log.Printf("min_num: %v", min_num)
	log.Printf("prestige_boost: %v", prestige_boost)
	log.Printf("MpC: %v", upgrade_stats.MpC)
	log.Printf("Mm: %v", upgrade_stats.Mm)
	
	session.Money += uint(math.Ceil(upgrade_stats.MpC * upgrade_stats.Mm * min_num * prestige_boost))
	log.Println("added value to money: ", uint(math.Ceil(upgrade_stats.MpC * upgrade_stats.Mm * min_num * prestige_boost)))

	session.Dishes -= uint(min_num)

	if session.LevelXP != 100 {
		session.LevelXP = math.Round(session.LevelXP + 1 + (session.PrestigeCurrent * 0.05) * 100) / 100
	}
	
	log.Println("dishes: ", session.Dishes)
	log.Println("money: ", session.Money)

	database.SaveSessionState(s.user_id, session)
	database.A.MarkChanged(s.user_id)

	return map[string]interface{}{
		"dishes": session.Dishes,
		"money": session.Money,
		"xp": session.LevelXP,
	}, SellRequest
}

func (s *SessionConn) ListUpgrades() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)

	return map[string]interface{}{
		"available": service.FilterUpgrades(session, false),
		"current": service.FilterUpgrades(session, true),
	}, ListRequest
}

func (s *SessionConn) LevelUp() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)
	next_level := database.LevelsXP[session.LevelRank + 1]

	if session.LevelRank == 100 {
		return map[string]interface{}{
			"current_rank": session.LevelRank,
			"current_xp": session.LevelXP,
		}, LevelUpRequest
	}

	if session.LevelXP == float64(next_level) {
		session.LevelRank += 1
		session.LevelXP = 0

		database.SaveSessionState(s.user_id, session)
		database.A.MarkChanged(s.user_id)

		return map[string]interface{}{
			"current_rank": session.LevelRank,
			"current_xp":   session.LevelXP,
			"next_xp":      database.LevelsXP[session.LevelRank + 1],
		}, LevelUpRequest
	}

	if session.LevelXP > float64(next_level) {
		session.LevelXP = math.Round((session.LevelXP - float64(next_level)) * 100) / 100
		session.LevelRank += 1

		database.SaveSessionState(s.user_id, session)
		database.A.MarkChanged(s.user_id)

		return map[string]interface{}{
			"current_rank": session.LevelRank,
			"current_xp":   session.LevelXP,
			"next_xp":      database.LevelsXP[session.LevelRank + 1],
		}, LevelUpRequest
	}

	database.SaveSessionState(s.user_id, session)
	database.A.MarkChanged(s.user_id)

	return map[string]interface{}{
		"current_rank": session.LevelRank,
		"current_xp":   session.LevelXP,
	}, LevelUpRequest
}

func (s *SessionConn) GetLevel() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)
	if session.LevelRank == 100 {
		return map[string]interface{}{
			"current_rank": session.LevelRank,
			"current_xp":   session.LevelXP,
		}, CheckLevelRequest
	}

	return map[string]interface{}{
		"current_rank": session.LevelRank,
		"current_xp":   session.LevelXP,
		"needed_xp":    database.LevelsXP[session.LevelRank + 1],
	}, CheckLevelRequest
}

func (s *SessionConn) ResetSession() (map[string]interface{}, RequestType) {
	session := database.GetSessionState(s.user_id)
	if session.PrestigeAccumulated < 1 {
		return map[string]interface{}{
			"message": "not enough prestige points",
		}, ErrorRequest
	}

	session.PrestigeCurrent += session.PrestigeAccumulated
	session.PrestigeBoost += math.Round((1 + 0.05 * session.PrestigeCurrent) * 10 ) / 10

	for i := range session.Upgrades {
		session.Upgrades[i] = 0
	}

	session.Money, session.Dishes, session.LevelRank, session.LevelXP, session.PrestigeAccumulated = 0, 0, 0, 0, 0

	database.SaveSessionState(s.user_id, session)
	database.A.MarkChanged(s.user_id)

	return map[string]interface{}{
		"message": "success",
	}, ResetRequest
}
