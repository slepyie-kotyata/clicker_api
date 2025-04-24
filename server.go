package main

import (
	"clicker_api/environment"
	"clicker_api/routes"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)




func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
	}))

	game := e.Group("/game")
	game.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(environment.GetVariable("ACCESS_TOKEN_SECRET")),
	}))

	refresh := e.Group("/refresh")
	refresh.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(environment.GetVariable("REFRESH_TOKEN_SECRET")),
	}))

	routes.InitEntryRoutes(e)
	routes.InitRefreshRoute(refresh)
	routes.InitGameRoutes(game)

	e.Start(":1323")

}

