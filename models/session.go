package models

type Session struct {
	ID         uint 		`gorm:"primary_key"`
	Money      uint
	Dishes     []Dish   	`gorm:"many2many:session_dishes;"`
	UserID     uint
}