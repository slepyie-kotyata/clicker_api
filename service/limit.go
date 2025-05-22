package service

import (
	"clicker_api/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func LimiterMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			if c.Request().Method == http.MethodGet {
				return next(c)
			}

			time_header := c.Request().Header.Get("X-timestamp")
			request_time := time.UnixMilli(int64(utils.StringToUint(time_header)))
			current_time := time.Now()

			diff := current_time.Sub(request_time)

			if diff >= time.Second || diff < 0 {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"status": 5,
					"message": "request blocked",
				})
			}

			return next(c)
		}
	}
}