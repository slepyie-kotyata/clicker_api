package custommiddleware

import (
	"clicker_api/pkg/format"
	"clicker_api/services/main_api/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {

            fmt.Printf("[JWT] %s %s -> middleware\n", c.Request().Method, c.Path())

            if c.Request().Method == http.MethodOptions {
                fmt.Println("[JWT] OPTIONS detected → skipping token check & returning 200 for CORS")
                return c.NoContent(200)
            }

            header := c.Request().Header.Get("Authorization")
            if header == "" {
                fmt.Println("[JWT] No Authorization header -> 401")
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "status":  1,
                    "message": "missing token",
                })
            }

            header_parts := strings.Split(header, " ")
            if len(header_parts) != 2 || header_parts[0] != "Bearer" {
                fmt.Printf("[JWT] Invalid Authorization format: %s\n", header)
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "status":  1,
                    "message": "invalid token format",
                })
            }

            token := header_parts[1]
            fmt.Println("[JWT] Token received, validating...")

            err := service.ValidateToken(token, secret)
            if err != nil {
                fmt.Printf("[JWT] Token invalid: %s\n", err.Error())
                return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                    "status":  1,
                    "message": err.Error(),
                })
            }

            uid := format.StringToUint(service.ExtractIDFromToken(token, secret))
            c.Set("id", uid)
            fmt.Printf("[JWT] Token valid. user_id=%d -> forwarding\n", uid)

            return next(c)
        }
    }
}