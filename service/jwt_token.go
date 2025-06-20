package service

import (
	"errors"
	"strings"
	"time"
	"clicker_api/secret"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(id string, is_access bool) string {
	var (
		token_secret string
		duration time.Duration
	)

	if is_access {
		token_secret = secret.Access_secret
		duration = time.Minute * 15
	} else {
		token_secret = secret.Refresh_secret
		duration = time.Hour * 168
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
	})

	signed_token, err :=  token.SignedString([]byte(token_secret))
	
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
