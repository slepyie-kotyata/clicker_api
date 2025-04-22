package models

type Session struct {
	ID         uint 		`json:"id" gorm:"primary_key"`
	Money      uint			`json:"money"`
	Dishes     uint         `json:"dishes"`
	UserID     uint			`json:"user_id"`
	Upgrades   []Upgrade    `gorm:"many2many:sesson_upgrades;"`
}