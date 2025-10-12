package routes

import (
	"clicker_api/ws"

	"github.com/labstack/echo/v4"
)

func InitWsRoutes(e *echo.Echo) {
	e.GET("/ws/", ws.ServeWs)
}