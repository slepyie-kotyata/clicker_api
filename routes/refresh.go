package routes

import (
	"clicker_api/handlers"
	"fmt"

	"github.com/labstack/echo/v4"
)

func InitRefreshRoute(refresh *echo.Group) {
    refresh.POST("", func(c echo.Context) error {
        fmt.Println("[REFRESH] POST /refresh request received")
        return handlers.RefreshTokens(c)
    })

    refresh.OPTIONS("", func(c echo.Context) error {
        fmt.Println("[REFRESH] OPTIONS preflight received -> return 200")
        return c.NoContent(200)
    })
}