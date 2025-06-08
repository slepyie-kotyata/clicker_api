package handlers

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"math"
	"net/http"

	"github.com/dariubs/percent"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BuyUpgrade(c echo.Context) error {
	user_id, _ := c.Get("id").(uint)
	upgrade_id := utils.StringToUint(c.Param("upgrade_id"))

	var (
		session      models.Session
		level_xp     models.LevelXP
		this_upgrade service.FilteredUpgrade
		exist        bool = false
		result_price uint = 0
		xp_increase  float64
	)

	database.DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", user_id).First(&session)

	if session.Level.Rank == 100 {
		xp_increase = 0
	} else {
		database.DB.Where("rank = ?", session.Level.Rank+1).Find(&level_xp)
		xp_increase = percent.Percent(5, int(level_xp.XP))

	}

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
		result_price = uint(math.Round(float64(this_upgrade.Price) * math.Pow(this_upgrade.PriceFactor, float64(this_upgrade.TimesBought))))
	}

	if session.Money < result_price {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough money, you have: " + utils.IntToString(int(session.Money)) + ", you need: " + utils.IntToString(int(result_price)),
		})
	}

	database.DB.Model(&session).Update("money", gorm.Expr("money - ?", result_price))
	database.DB.Model(&models.Level{}).Where("session_id = ?", session.ID).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", math.Abs(xp_increase)))
	database.DB.Model(&models.SessionUpgrade{}).Where("session_id = ? AND upgrade_id = ?", session.ID, upgrade_id).Select("times_bought").Updates(models.SessionUpgrade{TimesBought: this_upgrade.TimesBought + 1})

	database.DB.Preload("Level").First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"money":  session.Money,
		"xp":     session.Level.XP,
	})
}

func GetUpgrades(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var session models.Session
	database.DB.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	filtered_upgrades := service.FilterUpgrades(session, false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "0",
		"upgrades": filtered_upgrades,
	})
}