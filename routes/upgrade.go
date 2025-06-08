package routes

import (
	"clicker_api/handlers"

	"github.com/labstack/echo/v4"
)

func InitUpgradeRoutes(game *echo.Group) {
	game.PATCH("/upgrades/:upgrade_id", handlers.BuyUpgrade)
	game.GET("/upgrades", handlers.GetUpgrades)
}