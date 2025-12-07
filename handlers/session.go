package handlers

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitGame(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var (
		session models.Session
		user models.User
	)
	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	database.DB.Select("email").First(&user, id)

	if session.ID > 0 {
		filtered_upgrades := service.FilterUpgrades(&session, true)
		session.UserEmail = user.Email

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
			"upgrades": filtered_upgrades,
		})
	}

	new_session := models.Session{
		Money: 0,
		Dishes: 0,
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

	new_session.UserEmail = user.Email

	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
		"upgrades": make([]service.FilteredUpgrade, 0),
	})
}

func CookClick(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var session models.Session

	database.DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(&session, true))
	service.SetDefaults(&upgrade_stats)
	
	if upgrade_stats.HasDish == false {
		return c.JSON(http.StatusForbidden, map[string]string{
			"status":  "5",
			"message": "can't perform action",
		})
	}

	database.DB.Model(&session).Update("dishes", gorm.Expr("dishes + ?", uint(math.Ceil((1 + upgrade_stats.DpC) * upgrade_stats.Dm))))
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 10))
	database.DB.Preload("Level").First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"xp": session.Level.XP,
	})
}

func SellClick(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var (
		session models.Session
	)
	
	database.DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	
	if session.Dishes <= 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough dishes",
		})
	}
	
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(&session, true))
	service.SetDefaults(&upgrade_stats) 

	min_num := min(upgrade_stats.SpS, float64(session.Dishes))
	
	prestige_boost := session.Prestige.AccumulatedValue
	if prestige_boost == 0 {
		prestige_boost = 1
	}

	database.DB.Model(&session).Updates(map[string]interface{}{
		"money": gorm.Expr("money + ?", uint(math.Ceil(upgrade_stats.MpC * upgrade_stats.Mm * min_num * prestige_boost))),
		"dishes": gorm.Expr("dishes - ?", min_num),
	})

	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 10))
	database.DB.Preload("Level").First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"money":  session.Money,
		"xp":     session.Level.XP,
	})
}

func SessionReset(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var session models.Session
	database.DB.Preload("Prestige").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.Prestige.CurrentValue < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "2",
			"message": "not enough prestige points",
		})
	}

	p_boost := math.Round((1 + 0.05 * session.Prestige.CurrentValue) * 10 ) / 10
	p_value := session.Prestige.CurrentValue

	database.DB.Model(&models.SessionUpgrade{}).Where("session_id = ?", session.ID).Select("times_bought").Updates(&models.SessionUpgrade{TimesBought: 0})
	database.DB.Model(&models.Prestige{}).Where("session_id = ?", session.ID).Select("current_value").Updates(models.Prestige{CurrentValue: 0})
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Updates(map[string]interface{}{"rank": 0, "xp": 0})
	
	database.DB.Model(&session).Updates(map[string]interface{}{
		"money": 0,
		"dishes": 0,
		"prestige_boost": gorm.Expr("prestige_boost + ?", p_boost),
		"prestige_value": gorm.Expr("prestige_value + ?", p_value),
	})
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"message": "success",
	})
}