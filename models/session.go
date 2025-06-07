package models

type Session struct {
	ID         		uint 		`json:"-" gorm:"primary_key"`
	Money      		uint		`json:"money" gorm:"check:money >= 0"`
	Dishes     		uint		`json:"dishes" gorm:"check:dishes >= 0"`
	PrestigeValue 	float64		`json:"prestige_value" gorm:"check:prestige_value >= 0"`
	PrestigeBoost   float64		`json:"-" gorm:"check:prestige_boost >= 0"`
	UserID     		uint		`json:"user_id"`
	Level      		*Level		`json:"level" gorm:"foreignKey:SessionID"`
	Prestige   		*Prestige   `json:"prestige" gorm:"foreignKey:SessionID"`
	Upgrades   		[]Upgrade   `json:"-" gorm:"many2many:session_upgrades;"`
	UserEmail 		string 		`gorm:"-" json:"user_email"`
}