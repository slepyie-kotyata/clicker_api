package database

import (
	"clicker_api/models"
	"clicker_api/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

//TODO: удаление записи после выхода клиента

var (
	ctx = context.Background()
	Upgrades = FetchUpdates()
	LevelsXP = FetchLevelsXP()
)

func FetchLevelsXP() map[uint]uint {
	var levels []models.LevelXP
	DB.Find(&levels)

	levels_xp := make(map[uint]uint)
	for _, l := range levels {
		levels_xp[l.Rank] = l.XP
	}

	return levels_xp
}

func FetchUpdates() *[]models.Upgrade {
	var upgrades []models.Upgrade
	DB.Preload("Boost").Find(&upgrades)

	return &upgrades
}

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
		Upgrades: make(map[uint]uint, len(*Upgrades)),
	}

	DB.Where("session_id = ?", s.ID).Find(&session_upgrade)
	for _, u := range session_upgrade {
		session.Upgrades[u.UpgradeID] = u.TimesBought
	}

	fmt.Printf("RClient: %v, ctx: %v\n", RClient, ctx)

	data, _ := json.Marshal(session)

	fmt.Println(string(data))

	err := RClient.Set(ctx, utils.IntToString(int(s.UserID)), data, 0).Err()
	if err != nil {
		panic(err)
	}

	log.Println("session cashed")

	return &session
}

func SaveSessionState(user_id uint, s *models.SessionState) {
	data, _ := json.Marshal(s)
	if err := RClient.Set(ctx, utils.IntToString(int(user_id)), data, 0).Err(); err != nil {
		panic(err)
	}
}

func GetSessionState(user_id uint) *models.SessionState {
	result, err := RClient.Get(ctx, utils.IntToString(int(user_id))).Result()
	if err == redis.Nil {
        return nil
    }
	
	if err != nil {
		panic(err)
	}

	var session models.SessionState
	_ = json.Unmarshal([]byte(result), &session)

	return &session
}

func SetTTL(user_id uint) {
	_, err := RClient.Expire(ctx, utils.IntToString(int(user_id)), 20 * time.Minute).Result()
	if err != nil {
		panic(err)
	}
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
		
		tx.Model(&session).Updates(map[string]interface{}{
    		"money":  s.Money,
    		"dishes": s.Dishes,
		})

		tx.Model(session.Level).Updates(map[string]interface{}{
    		"rank": s.LevelRank,
    		"xp":   s.LevelXP,
		})

		tx.Model(session.Prestige).Updates(map[string]interface{}{
    		"current_value":         s.PrestigeCurrent,
    		"current_boost_value":   s.PrestigeBoost,
    		"accumulated_value":     s.PrestigeAccumulated,
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

//TODO: обновление списка SessionUpdates

func InitSession(id uint) *models.Session {
	var (
		session models.Session
		user models.User
	)

	log.Printf("loading...")
	
	DB.Preload("Prestige").Preload("Level").Where("user_id = ?", id).First(&session)
	DB.Select("email").First(&user, id)
	
	if session.ID > 0 {
		session.UserEmail = user.Email
		
		log.Printf("done!")
		return &session
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
	
	DB.Preload("Prestige").Preload("Level").Where("user_id = ?", id).First(&new_session)

	log.Printf("done!")
	
	return &new_session
}