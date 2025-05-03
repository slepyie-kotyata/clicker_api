package routes

import (
	passive_ws "clicker_api/handlers/passive_ws"
	"github.com/labstack/echo/v4"
)

func InitPassiveWS(e *echo.Echo) {
	e.GET("/passive/", passive_ws.ServeWS)
}