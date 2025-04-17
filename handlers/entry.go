package handlers

import (
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

//TODO: jwt auth, refresh handler, password hashing\decoding

func Authentication(c echo.Context) error {
	email, password := c.FormValue("email"), c.FormValue("password")

	var user models.User
	db.Preload("Password").Where("email = ? ", email).First(&user)
	if user.ID == 0 || !service.DoPasswordsMatch(user.Password.Hash, password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status": "1",
            "message": "invalid user info",
		})
	}
	
	access_token := service.NewToken(utils.IntToString(int(user.ID)), true)
	refresh_token := service.NewToken(utils.IntToString(int(user.ID)), false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"tokens": map[string]string {
            "access_token": access_token,
            "refresh_token": refresh_token,
        },
	})
}

func Registrate(c echo.Context) error {
	email, password  := c.FormValue("email"), c.FormValue("password")
	if (email == "" || password == "") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "2",
            "message": "not enough data",
		})
	}
			
	var user models.User	
	db.Where("email = ?", email).First(&user)	
	if (user.ID > 0) {
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "3",
            "message": "this user already exist",
		})
	}

	username:= strings.Split(email, "@")
	
	new_user := models.User {
		Username: username[0], 
		Email: email, 
		Password: models.Password {
			Hash: service.HashPassword(password),
		},
	}

	db.Create(&new_user)
	
	access_token := service.NewToken(utils.IntToString(int(new_user.ID)), true)
	refresh_token := service.NewToken(utils.IntToString(int(new_user.ID)), false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"tokens": map[string]string {
            "access_token": access_token,
            "refresh_token": refresh_token,
        },
	})
}