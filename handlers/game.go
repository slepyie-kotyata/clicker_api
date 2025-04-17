package handlers

import (
	"clicker_api/environment"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitGame(c echo.Context) error {
	id := service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), environment.GetVariable("ACCESS_TOKEN_SECRET"))
	new_id := utils.StringToUint(id)

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", new_id).First(&session)
	if session.ID > 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
		})
	}

	new_session := models.Session{
		Money: 0,
		UserID: new_id,
	}
	db.Create(&new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
	})
}