package database

import (
	"clicker_api/models"
	"log"
)

func InitSession(id uint) *models.Session {
	var (
		session models.Session
		user models.User
	)

	log.Printf("loading...")
	
	DB.Preload("Prestige").Preload("Level").Where("user_id = ?", id).First(&session)
	DB.Select("email").First(&user, id)
	
	if session.ID > 0 {
		log.Printf("done!")
		return &session
	}
	
	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		UserID: id,
		UserEmail: user.Email,
		Level: &models.Level{},
		Prestige: &models.Prestige{},
	}
	DB.Create(&new_session)
	

	var upgrades []models.Upgrade
	DB.Find(&upgrades)
	
	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		DB.Create(&session_upgrade)
	}
	
	DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	log.Printf("done!")
	
	return &new_session
}