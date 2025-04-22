package models

type SessionUpgrade struct {
	SessionID   uint	`gorm:"primary_key"` 
	UpgradeID   uint	`gorm:"primary_key"` 
	TimesBought uint
}