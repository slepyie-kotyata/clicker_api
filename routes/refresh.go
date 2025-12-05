package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitRefreshRoute(refresh *echo.Group) {
	refresh.POST("/", handlers.RefreshTokens)
	refresh.OPTIONS("/", func(c echo.Context) error {
    	return c.NoContent(200)
	})
}