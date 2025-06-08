package main

import (
	"clicker_api/custom_middleware"
	"clicker_api/routes"
	"clicker_api/secret"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)




func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://clicker.enjine.ru", "http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
	}))

	game := e.Group("/game", custommiddleware.LimiterMiddleware())
	game.Use(custommiddleware.JWTMiddleware(secret.Access_secret))
	
	refresh := e.Group("/refresh", custommiddleware.JWTMiddleware(secret.Refresh_secret))

	routes.InitEntryRoutes(e)
	routes.InitRefreshRoute(refresh)
	routes.InitSessionRoutes(game)
	routes.InitUpgradeRoutes(game)
	routes.InitLevelRoutes(game)
	routes.InitPassiveWS(e)

	e.Start(":1323")

}

