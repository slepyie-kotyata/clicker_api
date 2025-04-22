package models

type BoostType string

type Boost struct {
	ID			uint	 	`json:"id" gorm:"primary_key"`
	BoostType	BoostType	`json:"boost_type" gorm:"type:boost_type"`
	Value       uint		`json:"value"`
	UpgradeID   uint		`json:"-"`
}