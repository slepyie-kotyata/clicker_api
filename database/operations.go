package database

import (
	"clicker_api/models"
	"clicker_api/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var ctx = context.Background()

func CreateSessionState(s *models.Session) *models.SessionState {
	log.Println("start cashing..")
	var session_upgrade []models.SessionUpgrade

	session := models.SessionState{
		ID: s.ID,
		Money: s.Money,
		Dishes: s.Dishes,
		LevelRank: s.Level.Rank,
		LevelXP: s.Level.XP,
		PrestigeCurrent: s.Prestige.CurrentValue,
		PrestigeBoost: s.Prestige.CurrentBoostValue,
		PrestigeAccumulated: s.Prestige.AccumulatedValue,
		Upgrades: make(map[uint]uint, len(s.Upgrades)),
	}

	DB.Where("session_id = ?", s.ID).Find(&session_upgrade)
	for _, u := range session_upgrade {
		session.Upgrades[u.UpgradeID] = u.TimesBought
	}

	fmt.Printf("RClient: %v, ctx: %v\n", RClient, ctx)

	data, _ := json.Marshal(session)

	fmt.Println(string(data))

	err := RClient.Set(ctx, utils.IntToString(int(s.UserID)), data, 0)
	if err != nil {
		panic(err)
	}

	log.Println("session cashed")

	return &session
}

func SaveSession(s *models.SessionState) {
	log.Println("saving session...")

	DB.Transaction(func(tx *gorm.DB) error {
		var (
			session models.Session
		)
		
		if err := tx.Preload("Level").Preload("Prestige").First(&session, s.ID).Error; err != nil {
			return err
		}
		
		tx.Model(&session).Updates(models.Session{
			Money:  s.Money,
			Dishes: s.Dishes,
		})
		
		tx.Model(session.Level).Updates(models.Level{
			Rank: s.LevelRank,
			XP:   s.LevelXP,
		})
		
		tx.Model(session.Prestige).Updates(models.Prestige{
			CurrentValue: s.PrestigeCurrent,
			CurrentBoostValue: s.PrestigeBoost,
			AccumulatedValue: s.PrestigeAccumulated,
		})
		
		return nil
	})
	
	var session_upgrade []models.SessionUpgrade
	DB.Find(&session_upgrade, s.ID)

	for i := range session_upgrade {
    	session_upgrade[i].TimesBought = s.Upgrades[session_upgrade[i].UpgradeID]
	}

	DB.Save(&session_upgrade)

	log.Println("done!")
}

func InitSession(id uint) models.Session {
	var (
		session models.Session
		user models.User
	)

	log.Printf("loading...")
	
	DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	DB.Select("email").First(&user, id)
	
	if session.ID > 0 {
		session.UserEmail = user.Email
		
		log.Printf("done!")
		return session
	}
	
	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		UserID: id,
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
	
	new_session.UserEmail = user.Email
	
	DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	log.Printf("done!")
	
	return new_session
}