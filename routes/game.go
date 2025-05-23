package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitGameRoutes(game *echo.Group) {
	game.GET("/init", handlers.InitGame)
	game.PATCH("/reset", handlers.SessionReset)
	game.PATCH("/cook", handlers.CookClick)
	game.PATCH("/sell", handlers.SellClick)
	game.PATCH("/buy/:upgrade_id", handlers.BuyUpgrade)
	game.GET("/upgrades", handlers.GetUpgrades)
	game.GET("/levels", handlers.GetLevel)
	game.PATCH("/levels", handlers.UpdateLevel)
}
