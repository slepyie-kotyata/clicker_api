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

func InitGame(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	if session.ID > 0 {
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

	var first_upgrade models.Upgrade
	db.Preload("Boost").Where("icon_name = ?", "first_dish").First(&first_upgrade)
	db.Model(&new_session).Association("Upgrades").Append(&first_upgrade)

	db.Preload("Upgrades.Boost").First(&new_session, new_session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
	})
}

func CookClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), environment.GetVariable("ACCESS_TOKEN_SECRET")))

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.Upgrades == nil {
		return c.JSON(http.StatusInsufficientStorage, map[string]string{
			"status": "4",
			"dishes": "no upgrades found",
		})
	}

	var (
		total_dishes_multiplier uint = 0;
		total_dishes_per_click uint = 0;
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

	db.Model(&session).Select("dishes").Updates(models.Session{Dishes: session.Dishes + (total_dishes_per_click * 5 * total_dishes_multiplier)})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
	})
}