package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitGameRoutes(access *echo.Group) {
	access.GET("/game/init", handlers.InitGame)
}