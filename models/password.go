package models

type Password struct {
	ID     uint		`gorm:"primary_key"`
	Hash   string
	UserID uint
}