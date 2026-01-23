package models

type UpgradeType string

type Upgrade struct {
	ID             uint		       `json:"id" gorm:"primary_key"`
	Name           string		   `json:"name"`
	IconName       string		   `json:"icon_name"`
	UpgradeType    UpgradeType     `json:"upgrade_type" gorm:"type:upgrade_type"`
	PriceFactor    float64		   `json:"price_factor"`
	Price          uint			   `json:"price"`
	AccessLevel    uint			   `json:"access_level"`
	Boost          Boost           `json:"boost" gorm:"foreignKey:UpgradeID"`
}