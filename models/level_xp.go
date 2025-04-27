package models

type LevelXP struct {
	ID   uint `gorm:"primary_key"`
	Rank uint
	XP   uint
}