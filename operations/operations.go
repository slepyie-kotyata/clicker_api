package operations

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"context"
	"encoding/json"
	"log"
)

type SessionResponse struct {
	UserID     		uint		`json:"user_id"`
	UserEmail 		string 		`json:"user_email"`
	Money      		uint		`json:"money"`
	Dishes     		uint		`json:"dishes"`
	Level      		struct {
		Rank		uint	`json:"rank"`	
		XP			float64	`json:"xp"`
	}	`json:"level"`
	Prestige   		struct {
		CurrentValue 		float64		`json:"current_value"`
		CurrentBoostValue 	float64		`json:"current_boost_value"`
		AccumulatedValue   	float64		`json:"accumulated_value"`
	}   `json:"prestige"`
	Upgrades   		struct {
		Available	[]service.FilteredUpgrade	`json:"available"`
		Current 	[]service.FilteredUpgrade	`json:"current"`
	}   `json:"upgrades"`
}

func NewSessionResponse(session *models.Session) SessionResponse {
	return SessionResponse{
			UserID: session.UserID,
			UserEmail: session.UserEmail,
			Money: session.Money,
			Dishes: session.Dishes,
			Level: struct {
				Rank uint    `json:"rank"`
				XP   float64 `json:"xp"`
			}{
				Rank: session.Level.Rank,
				XP:   session.Level.XP,
			},
			Prestige: struct {
				CurrentValue       float64 `json:"current_value"`
				CurrentBoostValue  float64 `json:"current_boost_value"`
				AccumulatedValue   float64 `json:"accumulated_value"`
			}{
				CurrentValue:      session.Prestige.CurrentValue,
				CurrentBoostValue: session.Prestige.CurrentBoostValue,
				AccumulatedValue:  session.Prestige.AccumulatedValue,
			},
			Upgrades: struct {
				Available []service.FilteredUpgrade `json:"available"`
				Current   []service.FilteredUpgrade `json:"current"`
			}{
				Available: service.FilterUpgrades(session, false),
				Current: service.FilterUpgrades(session, true),
			},
		}
}

var ctx = context.Background()

func CreateSessionState(s *models.Session) *models.SessionState {
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

	database.DB.Where("session_id = ?", s.ID).Find(&session_upgrade)
	for _, u := range session_upgrade {
		session.Upgrades[u.UpgradeID] = u.TimesBought
	}

	data, _ := json.Marshal(session)
	err := database.RClient.Set(ctx, utils.IntToString(int(s.UserID)), data, 0)
	if err != nil {
		panic(err)
	}

	return &session
}

func InitSession(id uint) models.Session {
	var (
		session models.Session
		user models.User
	)

	log.Printf("loading...")
	
	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	database.DB.Select("email").First(&user, id)
	
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
	
	database.DB.Create(&new_session)
	
	var upgrades []models.Upgrade
	database.DB.Find(&upgrades)
	
	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		database.DB.Create(&session_upgrade)
	}
	
	new_session.UserEmail = user.Email
	
	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	log.Printf("done!")
	
	return new_session
}