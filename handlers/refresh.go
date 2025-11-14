package handlers

import (
	"clicker_api/environment"
	"clicker_api/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RefreshTokens(c echo.Context) error {	
	id := service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), environment.GetVariable("REFRESH_TOKEN_SECRET"))

 	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"tokens": map[string]string{
			"access_token": service.NewToken(id, true),
			"refresh_token": service.NewToken(id, false),
		},
	})
}