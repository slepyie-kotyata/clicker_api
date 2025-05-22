package service

import (
	"clicker_api/environment"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewToken(id string, is_access bool) string {
	var (
		secret string
		duration time.Duration
	)

	if is_access {
		secret = environment.GetVariable("ACCESS_TOKEN_SECRET")
		duration = time.Minute * 15
	} else {
		secret = environment.GetVariable("REFRESH_TOKEN_SECRET")
		duration = time.Hour * 168
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
	})

	signed_token, err :=  token.SignedString([]byte(secret))
	
	if err != nil {
		return ""
	}

	return signed_token
}

func ParseToken(token string, secret string) (*jwt.Token, error) {
	parsed_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return parsed_token, err
}

func ExtractIDFromToken(full_token string, secret string) string {
	token_parts := strings.Split(full_token, " ")

	var token *jwt.Token

	if len(token_parts) > 1 {
		token, _ = ParseToken(token_parts[1], secret)
	} else {
		token, _ = ParseToken(token_parts[0], secret)
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	id := claims["id"].(string)

	return id
}

func ValidateToken(token_string string, secret string) error {
	token, err := ParseToken(token_string, secret)
	if err != nil || token == nil || !token.Valid {
		return errors.New("token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	id, ok := claims["id"]
	
	if !ok {
		return errors.New("invalid token claims")
	}
	
	if _, ok := id.(string); !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if ok {
		if int64(exp) < time.Now().Unix() {
			return errors.New("token expired")
		}
	} else {
		return errors.New("invalid or missing expiration time")
	}

	return nil
}

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func (c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": 1,
					"message": "missing token",
				})
			}

			header_parts := strings.Split(header, " ")
			if len(header_parts) != 2 || header_parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": 1,
					"message": "invalid token format",
				})
			}

			token := header_parts[1]

			err := ValidateToken(token, secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": 1,
					"message": err.Error(),
				})
			}

			c.Set("id", ExtractIDFromToken(token, secret))

			return next(c)
		}
	}
}
