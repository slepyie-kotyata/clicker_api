package service

import (
	"clicker_api/environment"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func ExtractIDFromToken(header string, secret string) string {
	header_parts := strings.Split(header, " ")
	token, _ := ParseToken(header_parts[1], secret)

	claims, _ := token.Claims.(jwt.MapClaims)

	id := claims["id"].(string)

	return id
}

