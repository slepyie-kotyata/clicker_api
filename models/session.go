package models

type Session struct {
	ID         uint 			`json:"-" gorm:"primary_key"`
	Money      uint				`json:"money"`
	Dishes     		uint		`json:"dishes"`
	PrestigeValue 	float64		`json:"prestige_value"`
	UserID     		uint		`json:"user_id"`
	Level      		*Level		`json:"level" gorm:"foreignKey:SessionID"`
	Prestige   		*Prestige   `json:"prestige" gorm:"foreignKey:SessionID"`
	Upgrades   		[]Upgrade   `json:"-" gorm:"many2many:session_upgrades;"`
}