package models

type User struct {
	ID       uint		`json:"id" gorm:"primary_key"`
	Email    string		`json:"email"`
	Password Password	`json:"-" gorm:"foreignKey:UserID"`
	Session  Session    `json:"-" gorm:"foreignKey:UserID"`
}
