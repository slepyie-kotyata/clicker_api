package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitSessionRoutes(game *echo.Group) {
	game.GET("/sessions", handlers.InitGame)
	game.PATCH("/sessions/reset", handlers.SessionReset)
	game.PATCH("/sessions/cook", handlers.CookClick)
	game.PATCH("/sessions/sell", handlers.SellClick)
}
