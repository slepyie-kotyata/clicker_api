package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitGameRoutes(game *echo.Group) {
	game.GET("/init", handlers.InitGame)
	// game.PATCH("/cook", handlers.CookClick)
}