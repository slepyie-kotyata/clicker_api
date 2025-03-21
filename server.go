package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primary_key"`

	Username string
	Email string
	Password Password
}

type Password struct {
	ID uint `gorm:"primary_key"`
	Hash string
	UserID uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	  panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Password{})
	if err != nil {
	  panic("failed to connect database")
	}

	e := echo.New()

	e.POST("/reg", func(c echo.Context) error {
		username, email, password  := c.FormValue("username"), c.FormValue("email"), c.FormValue("password") // Получаем имя из URL
		if username == " " || email == " " || password == " " {
			return c.String(http.StatusBadRequest, "Недостаточно данных")
		}
			
		var user User	
		db.Where("username = ? OR email = ?", username, email).First(&user)	
		if (user.ID !=0){
			return c.String(http.StatusConflict, "Пользователь уже существует")
		}
		new_user := User{
			Username: username, 
			Email: email, 
			Password: Password{
				Hash: password,
			},
		}
		db.Create(&new_user)

		return c.JSON(http.StatusOK, &new_user)
	})

	e.POST("/auth", func(c echo.Context) error {
		email, password  :=  c.FormValue("email"), c.FormValue("password") // Получаем имя из URL
		if email == " " || password == " " {
			return c.String(http.StatusBadRequest, "Недостаточно данных")
		}
			
		var user User	
		db.Preload("Password").Where("email = ? ", email).First(&user)	
		if (user.ID == 0 || user.Password.Hash != password){
			return c.String(http.StatusUnauthorized, "Такого пользователя нет")
		}

		return c.JSON(http.StatusOK, &user)
	})

	e.Start(":1323")

}

