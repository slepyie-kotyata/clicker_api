package models

type Level struct {
	ID			uint	`json:"-" gorm:"primary_key"` 
	Rank		uint	`json:"rank"`	
	XP			uint	`json:"xp"`
	SessionID 	uint	`json:"-"`	
}