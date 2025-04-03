package handlers

import (
	"clicker_api/environment"
	"clicker_api/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RefreshTokens(c echo.Context) error {
	header := strings.Split(c.Request().Header.Get("Authorization"), " ")
	id := service.ExtractIDFromToken(header[1], environment.GetVariable("REFRESH_TOKEN_SECRET"))

 	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "0",
		"tokens": map[string]string{
			"access_token": service.NewToken(id, true),
			"refresh_token": service.NewToken(id, false),
		},
	})
}