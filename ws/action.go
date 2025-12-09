package ws

import (
	"clicker_api/database"
	"clicker_api/service"
	"math"
)

// func (s *SessionConn) Buy() map[string]interface{} {
// 	s.session = database.GetSessionState(s.user_id)
// }

func (s *SessionConn) Cook() map[string]interface{} {
	s.session = database.GetSessionState(s.user_id)

	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(s.session, true))
	service.SetDefaults(&upgrade_stats)
	
	if !upgrade_stats.HasDish {
		return map[string]interface{}{
			"message": "can't perform this action",
		}
	}

	s.session.Dishes += uint(math.Ceil((1 + upgrade_stats.DpC) * upgrade_stats.Dm))
	s.session.LevelXP = math.Round((s.session.LevelXP + 10) * 100) / 100

	database.SaveSessionState(s.user_id, s.session)

	return map[string]interface{}{
		"dishes": s.session.Dishes,
		"xp": s.session.LevelXP,
	}
}

