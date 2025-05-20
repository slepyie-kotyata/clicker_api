package main

import (
	"clicker_api/handlers"
	"clicker_api/routes"
	"clicker_api/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)




func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
	}))

	game := e.Group("/game", service.JWTMiddleware(handlers.Access_secret))
	refresh := e.Group("/refresh", service.JWTMiddleware(handlers.Refresh_secret))

	routes.InitEntryRoutes(e)
	routes.InitRefreshRoute(refresh)
	routes.InitGameRoutes(game)
	routes.InitPassiveWS(e)

	e.Start(":1323")

}

