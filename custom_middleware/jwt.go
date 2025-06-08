package custommiddleware

import (
	"clicker_api/service"
	"clicker_api/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  1,
					"message": "missing token",
				})
			}

			header_parts := strings.Split(header, " ")
			if len(header_parts) != 2 || header_parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  1,
					"message": "invalid token format",
				})
			}

			token := header_parts[1]

			err := service.ValidateToken(token, secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status":  1,
					"message": err.Error(),
				})
			}

			c.Set("id", utils.StringToUint(service.ExtractIDFromToken(token, secret)))

			return next(c)
		}
	}
}