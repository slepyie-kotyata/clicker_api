package service

import (
	"clicker_api/utils"
	"fmt"
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
			fmt.Println("time_header ", time_header)

			request_time := time.UnixMilli(int64(utils.StringToUint(time_header)))
			fmt.Println("request_time formatted ", request_time)

			current_time := time.Now()
			fmt.Println("current_time ", current_time)

			diff := current_time.Sub(request_time)

			fmt.Println("diff ", diff)

			if diff >= time.Second * 2 || diff < 0 {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"status": 5,
					"message": "request blocked",
				})
			}

			return next(c)
		}
	}
}