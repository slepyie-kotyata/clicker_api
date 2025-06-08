package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitSessionRoutes(game *echo.Group) {
	game.GET("/sessions", handlers.InitGame)
	game.PATCH("/session/reset", handlers.SessionReset)
	game.PATCH("/session/cook", handlers.CookClick)
	game.PATCH("/session/sell", handlers.SellClick)
}
