package main

import (
	"clicker_api/custom_middleware"
	"clicker_api/routes"
	"clicker_api/secret"
	"clicker_api/ws"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    	AllowOrigins: []string{
      		"http://localhost:4200",
      		"https://clicker.enjine.ru",
      		"https://enjine.ru",
    	},
    	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
    	AllowHeaders: []string{
      		echo.HeaderOrigin,
      		echo.HeaderContentType,
      		echo.HeaderAccept,
      		echo.HeaderAuthorization,
    	},
    	AllowCredentials: true,
  }))

	game := e.Group("/game", custommiddleware.LimiterMiddleware())
	game.Use(custommiddleware.JWTMiddleware(secret.Access_secret))
	
	refresh := e.Group("/refresh")
	refresh.Use(custommiddleware.JWTMiddleware(secret.Refresh_secret))
	

	routes.InitEntryRoutes(e)
	routes.InitRefreshRoute(refresh)
	routes.InitSessionRoutes(game)
	routes.InitUpgradeRoutes(game)
	routes.InitLevelRoutes(game)
	routes.InitPassiveWS(e)
	routes.InitWsRoutes(e)

	e.Start(":1323")

}

