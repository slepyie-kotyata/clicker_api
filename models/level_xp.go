package models

import (
	"math"

	"gorm.io/gorm"
)

type LevelXP struct {
	ID   uint `gorm:"primary_key"`
	Rank uint
	XP   uint
}

func (l *LevelXP) BeforeCreate(tx *gorm.DB) (err error) {
	l.XP = uint(math.Floor(10 * math.Pow(float64(l.Rank), 1.5)))
	return
}