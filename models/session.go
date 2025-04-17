package models

type Session struct {
	ID         uint 		`json:"id" gorm:"primary_key"`
	Money      uint			`json:"money"`
	Upgrades   []Upgrade   	`json:"upgrades" gorm:"many2many:session_upgrades;"`
	UserID     uint			`json:"user_id"`
}