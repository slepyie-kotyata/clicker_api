package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitLevelRoutes(game *echo.Group) {
	game.GET("/levels", handlers.GetLevel)
	game.PATCH("/levels", handlers.UpdateLevel)
}