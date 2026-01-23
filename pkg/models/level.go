package models

type Level struct {
	ID			uint	`json:"-" gorm:"primary_key"` 
	Rank		uint	`json:"rank" gorm:"default:0;check:rank >= 0"`	
	XP			float64	`json:"xp" gorm:"default:0;check:xp >= 0.0"`
	SessionID 	uint	`json:"-"`	
}