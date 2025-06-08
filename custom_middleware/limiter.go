package custommiddleware

import (
	"clicker_api/utils"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var debug = false

func debugLog(format string, a ...interface{}) {
	if debug {
		log.Printf(format, a...)
	}
}

func LimiterMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			if c.Request().Method == http.MethodGet {
				return next(c)
			}

			time_header := c.Request().Header.Get("X-timestamp")
			debugLog("time_header: %s", time_header)

			request_time := time.UnixMilli(int64(utils.StringToUint(time_header)))
			debugLog("request_time: %v", request_time)

			current_time := time.Now()
			debugLog("current_time: %v", current_time)

			diff := current_time.Sub(request_time)
			debugLog("diff: %v", diff)

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