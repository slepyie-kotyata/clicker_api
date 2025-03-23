package service

import (
	"clicker_api/environment"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(id uint, is_access bool) string {
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


