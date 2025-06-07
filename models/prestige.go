package models

type Prestige struct {
	ID				uint		`json:"-" gorm:"primary_key"`
	CurrentValue 	float64		`json:"current_value" gorm:"default:0;check:current_value >= 0.0"`
	SessionID       uint		`json:"-"`
}