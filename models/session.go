package models

type Session struct {
	ID         uint 		`json:"id" gorm:"primary_key"`
	Money      uint			`json:"money"`
	Dishes     uint         `json:"dishes"`
	UserID     uint			`json:"user_id"`
	Level      Level        `json:"level" gorm:"foreignKey:SessionID"`
	Upgrades   []Upgrade    `json:"-" gorm:"many2many:session_upgrades;"`
}