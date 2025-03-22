package routes

import (
	"clicker_api/handlers"
	"github.com/labstack/echo/v4"
)

func InitEntryRoutes(e *echo.Echo) {
	e.POST("/reg", handlers.Registrate)
	e.POST("/auth", handlers.Authentication)
}