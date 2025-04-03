package models

type Dish struct {
	ID             uint		`gorm:"primary_key"`
	Name           string
	Price          uint
	MoneyPerClick  uint
}