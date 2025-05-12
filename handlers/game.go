package handlers

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)





func InitGame(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session
	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.ID > 0 {
		filtered_upgrades := service.FilterUpgrades(session, true)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
			"upgrades": filtered_upgrades,
		})
	}

	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		PrestigeValue: 0,
		UserID: id,
		Level: &models.Level{},
		Prestige: &models.Prestige{},
	}
	database.DB.Create(&new_session)

	var upgrades []models.Upgrade
	database.DB.Find(&upgrades)

	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		database.DB.Create(&session_upgrade)
	}

	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
		"upgrades": make([]service.FilteredUpgrade, 0),
	})
}

func CookClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session

	database.DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(session, true))
	
	if upgrade_stats.HasDish == false {
		return c.JSON(http.StatusForbidden, map[string]string{
			"status":  "5",
			"message": "can't perform action",
		})
	}

	database.DB.Model(&session).Update("dishes", gorm.Expr("dishes + ?", uint(math.Ceil((1 + upgrade_stats.DpC) * upgrade_stats.Dm))))
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 0.2))
	database.DB.Preload("Level").First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"xp": session.Level.XP,
	})
}

func SellClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var (
		session models.Session
	)
	
	database.DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(session, true))
	
	min_num := min(upgrade_stats.SpS, float64(session.Dishes))
	
	if session.Dishes <= 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough dishes",
		})
	}

	database.DB.Model(&session).Updates(map[string]interface{}{
		"money": gorm.Expr("money + ?", uint(math.Ceil(upgrade_stats.MpC * upgrade_stats.Mm * min_num))),
		"dishes": gorm.Expr("dishes - ?", min_num),
	})
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 0.2))
	database.DB.Preload("Level").First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"money":  session.Money,
		"xp":     session.Level.XP,
	})
}

func BuyUpgrade(c echo.Context) error {
	user_id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	upgrade_id := utils.StringToUint(c.Param("upgrade_id"))

	var (
		session      models.Session
		this_upgrade service.FilteredUpgrade
		exist        bool = false
		result_price uint = 0
	)

	database.DB.Preload("Upgrades.Boost").Where("user_id = ?", user_id).First(&session)

	for _, upgrade := range service.FilterUpgrades(session, false) {
		if upgrade.ID == upgrade_id {
			this_upgrade = upgrade
			exist = true
		}
	}

	if !exist {
		return c.JSON(http.StatusNotFound, map[string]string{
			"status":  "4",
			"message": "upgrade not found",
		})
	}

	if this_upgrade.TimesBought == 0 {
		result_price = this_upgrade.Price
	} else {
		result_price = uint(math.Ceil(float64(this_upgrade.Price) * this_upgrade.PriceFactor * float64(this_upgrade.TimesBought)))
	}

	if session.Money < result_price {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough money, you have: " + utils.IntToString(int(session.Money)) + ", you need: " + utils.IntToString(int(result_price)),
		})
	}

	database.DB.Model(&session).Update("money", gorm.Expr("money - ?", result_price))
	database.DB.Model(&models.SessionUpgrade{}).Where("session_id = ? AND upgrade_id = ?", session.ID, upgrade_id).Select("times_bought").Updates(models.SessionUpgrade{TimesBought: this_upgrade.TimesBought + 1})

	database.DB.First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"money":  session.Money,
	})
}

func GetUpgrades(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session
	database.DB.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	filtered_upgrades := service.FilterUpgrades(session, false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "0",
		"upgrades": filtered_upgrades,
	})
}

func GetLevel(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	var (
		session models.Session
		level   models.LevelXP
	)

	database.DB.Preload("Level").Where("user_id = ?", id).First(&session)
	database.DB.Where("rank = ?", session.Level.Rank + 1).First(&level)

	if level.ID == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": 0,
			"current_rank": session.Level.Rank,
			"current_xp": session.Level.XP,
			"needed_xp": session.Level.XP,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"current_rank": session.Level.Rank,
		"current_xp": session.Level.XP,
		"needed_xp": level.XP,
	})
}

func UpdateLevel(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	var (
		level   models.Level
		next_level models.LevelXP
	)

	database.DB.Where("session_id = (?)", database.DB.Model(&models.Session{}).Select("id").Where("user_id = ?", id),).First(&level)
	database.DB.Where("rank = ?", level.Rank + 1).First(&next_level)
	
	if next_level.ID == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp":   level.XP,
		})
	}

	if level.XP == float64(next_level.XP){
		database.DB.Model(&level).Updates(map[string]interface{}{
			"xp": 0,
			"rank": gorm.Expr("rank + ?", 1),
		})
		database.DB.First(&level, level.ID)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp": level.XP,
		})
	}

	if level.XP > float64(next_level.XP) {
		database.DB.Model(&level).Updates(map[string]interface{}{
			"xp": gorm.Expr("ROUND(xp - ?, 2)", next_level.XP),
			"rank": gorm.Expr("rank + ?", 1),
		})
		database.DB.First(&level, level.ID)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp": level.XP,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"current_rank": level.Rank,
		"current_xp": level.XP,
	})
}

func SessionReset(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session
	database.DB.Preload("Prestige").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.Prestige.CurrentValue < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "2",
			"message": "not enough prestige points",
		})
	}

	b := math.Round((1 + 0.05 * session.Prestige.CurrentValue) * 10 ) / 10

	database.DB.Model(&models.SessionUpgrade{}).Where("session_id = ?", session.ID).Select("times_bought").Updates(&models.SessionUpgrade{TimesBought: 0})
	database.DB.Model(&models.Prestige{}).Where("session_id = ?", session.ID).Select("current_value").Updates(models.Prestige{CurrentValue: 0})
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Updates(map[string]interface{}{"rank": 0, "xp": 0})
	
	database.DB.Model(&session).Updates(map[string]interface{}{
		"money": 0,
		"dishes": 0,
		"prestige_value": gorm.Expr("prestige_value + ?", b),
	})
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"message": "success",
	})
}