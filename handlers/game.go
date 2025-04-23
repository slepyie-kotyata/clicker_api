package handlers

import (
	"clicker_api/environment"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

var secret string = environment.GetVariable("ACCESS_TOKEN_SECRET")

func filterUpgrades(session models.Session) []models.Upgrade {
	filtered_upgrades := make([]models.Upgrade, 0)

	for _, upgrade := range session.Upgrades {
		var session_upgrade models.SessionUpgrade
		db.Where("session_id = ? AND upgrade_id = ?", session.ID, upgrade.ID).First(&session_upgrade)
		if session_upgrade.TimesBought > 0 {
			filtered_upgrades = append(filtered_upgrades, upgrade)
		}
	}

	return filtered_upgrades 
} 

func InitGame(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.ID > 0 {
		session.Upgrades = filterUpgrades(session)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
		})
	}

	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		UserID: id,
	}

	db.Create(&new_session)
	var upgrades []models.Upgrade
	db.Find(&upgrades)

	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		db.Create(&session_upgrade)
	}

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)
	new_session.Upgrades = filterUpgrades(new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
	})
}

func CookClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), environment.GetVariable("ACCESS_TOKEN_SECRET")))

	var session models.Session

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	session.Upgrades = filterUpgrades(session)

	var (
		total_dishes_multiplier float32 = 0;
		total_dishes_per_click float32 = 0;
	)

	for _, upgrade := range session.Upgrades {
		if upgrade.Boost.BoostType == "dM" {
			total_dishes_multiplier += upgrade.Boost.Value
		}

		if upgrade.Boost.BoostType == "dPc" {
			total_dishes_per_click += upgrade.Boost.Value
		}
	}

	if total_dishes_multiplier == 0 {
		total_dishes_multiplier = 1
	}

	db.Model(&session).Select("dishes").Updates(models.Session{Dishes: session.Dishes + uint((1 + total_dishes_per_click) * 5 * total_dishes_multiplier)})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
	})
}

func SellClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), environment.GetVariable("ACCESS_TOKEN_SECRET")))

	var session models.Session

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	session.Upgrades = filterUpgrades(session)

	var (
		total_money_multiplier float32 = 0;
		total_money_per_click float32 = 0;
	)

	for _, upgrade := range session.Upgrades {
		if upgrade.Boost.BoostType == "mM" {
			total_money_multiplier += upgrade.Boost.Value
		}

		if upgrade.Boost.BoostType == "mPc" {
			total_money_per_click += upgrade.Boost.Value
		}
	}

	if total_money_multiplier == 0 {
		total_money_multiplier = 1
	}

	if session.Dishes <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "2",
			"message": "zero dishes",
		})
	}

	db.Model(&session).Select("dishes", "money").Updates(models.Session{Dishes: session.Dishes - 1, Money: session.Money + uint((total_money_per_click) * 5 * total_money_multiplier)})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"money": session.Money,
	})
}