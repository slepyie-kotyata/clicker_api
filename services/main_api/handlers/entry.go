package handlers

import (
	"clicker_api/pkg/format"
	"clicker_api/pkg/models"
	"clicker_api/services/main_api/database"
	"clicker_api/services/main_api/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authentication(c echo.Context) error {
	email, password := c.FormValue("email"), c.FormValue("password")

	var user models.User
	database.DB.Preload("Password").Where("email = ? ", email).First(&user)
	if user.ID == 0 || !service.DoPasswordsMatch(user.Password.Hash, password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status": "1",
            "message": "invalid user info",
		})
	}
	
	access_token := service.NewToken(format.IntToString(int(user.ID)), true)
	refresh_token := service.NewToken(format.IntToString(int(user.ID)), false)

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
	database.DB.Where("email = ?", email).First(&user)	
	if (user.ID > 0) {
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "3",
            "message": "this user already exist",
		})
	}
	
	new_user := models.User {
		Email: email, 
		Password: models.Password {
			Hash: service.HashPassword(password),
		},
	}

	database.DB.Create(&new_user)
	
	access_token := service.NewToken(format.IntToString(int(new_user.ID)), true)
	refresh_token := service.NewToken(format.IntToString(int(new_user.ID)), false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"tokens": map[string]string {
            "access_token": access_token,
            "refresh_token": refresh_token,
        },
	})
}