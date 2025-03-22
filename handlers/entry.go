package handlers

import (
	"clicker_api/models"
	"net/http"
	"github.com/labstack/echo/v4"
)

func Authentication(c echo.Context) error {
	email, password := c.FormValue("email"), c.FormValue("password")
	if email == " " || password == " " {
		return c.String(http.StatusBadRequest, "Недостаточно данных")
	}

	var user models.User
	db.Preload("Password").Where("email = ? ", email).First(&user)
	if user.ID == 0 || user.Password.Hash != password {
		return c.String(http.StatusUnauthorized, "Такого пользователя нет")
	}

	return c.JSON(http.StatusOK, &user)
}

func Registrate(c echo.Context) error {
	username, email, password  := c.FormValue("username"), c.FormValue("email"), c.FormValue("password")
	if username == " " || email == " " || password == " " {
		return c.String(http.StatusBadRequest, "Недостаточно данных")
	}
			
	var user models.User	
	db.Where("username = ? OR email = ?", username, email).First(&user)	
	if (user.ID !=0) {
		return c.String(http.StatusConflict, "Пользователь уже существует")
	}
	
	new_user := models.User {
		Username: username, 
		Email: email, 
		Password: models.Password {
			Hash: password,
		},
	}

	db.Create(&new_user)

	return c.JSON(http.StatusOK, &new_user)
}