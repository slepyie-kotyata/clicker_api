package models

type Level struct {
	ID			uint	`json:"-" gorm:"primary_key"` 
	Rank		uint	`json:"rank" gorm:"default:0"`	
	XP			uint	`json:"xp" gorm:"default:0"`
	SessionID 	uint	`json:"-"`	
}