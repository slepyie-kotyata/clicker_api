package passivews

import (
	"clicker_api/handlers"
	"clicker_api/models"
	"fmt"
	"math"
	"time"

	"github.com/gorilla/websocket"
)

type SessionMessage struct {
	Money  			uint 		`json:"money"`
	Dishes 			uint 		`json:"dishes"`
	Rank   			uint 		`json:"rank"`
	XP     			float64  	`json:"xp"`
	PrestigeCurrent float64 	`json:"prestige_current"`
}

type Session struct {
	Session  models.Session
	Client   *websocket.Conn
	Messages chan SessionMessage
}

func (s *Session) UpdateSessionState(seconds uint) {
	filtered_upgrades := handlers.FilterUpgrades(s.Session, true)
	var (
		total_money_per_second float64 = 0
		total_dishes_per_second float64 = 0
		total_money_passive_multiplier float64 = 0
		total_dishes_passive_multiplier float64 = 0
	)

	dish_exist := false

	for _, upgrade := range filtered_upgrades {
		if upgrade.UpgradeType == "dish" && dish_exist == false {
			dish_exist = true
		}
		if upgrade.Boost.BoostType == "mPs" {
			total_money_per_second += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}
		if upgrade.Boost.BoostType == "dPs" {
			total_dishes_per_second += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}
		if upgrade.Boost.BoostType == "mpM" {
			total_money_passive_multiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}
		if upgrade.Boost.BoostType == "dpM" {
			total_dishes_passive_multiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}
	}

	if dish_exist == false {
		return
	}

	if total_money_passive_multiplier == 0 {
		total_money_passive_multiplier = 1

		if total_dishes_passive_multiplier == 0 {
			total_dishes_passive_multiplier = 1
		}
	}

	handlers.DB.Model(&s.Session).Select("dishes").Updates(models.Session{Dishes: s.Session.Dishes + uint(math.Ceil((1 + total_dishes_per_second) * total_dishes_passive_multiplier))})

	if s.Session.Dishes <= 0 {
		return
	}

	handlers.DB.Model(&s.Session).Select("dishes", "money").Updates(models.Session{Dishes: s.Session.Dishes - 1, Money: s.Session.Money + uint(math.Ceil((total_money_per_second) * total_money_passive_multiplier))})
	new_xp := math.Round((s.Session.Level.XP + 0.2 ) * 100) / 100
	handlers.DB.Model(&models.Level{}).Where("session_id = ?", s.Session.ID).Select("xp").Updates(map[string]interface{}{"xp": new_xp})

	s.Messages <- SessionMessage{
		Money: s.Session.Money,
		Dishes: s.Session.Dishes,
		Rank: s.Session.Level.Rank,
		XP: s.Session.Level.XP,
		PrestigeCurrent: s.Session.Prestige.CurrentValue,
	}
}

func (s *Session) StartPassiveLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		fmt.Println("stopped")
		ticker.Stop()
	}()

	for {
		select {
		case <- ticker.C:
			s.UpdateSessionState(3)
		}
	}
}

func (s *Session) HandleConnection(sm *SessionManager) {
	defer sm.CloseSession(s.Session.ID)

	session_message := map[string]interface{}{"message":SessionMessage{
		Money: s.Session.Money,
		Dishes: s.Session.Dishes,
		Rank: s.Session.Level.Rank,
		XP: s.Session.Level.XP,
		PrestigeCurrent: s.Session.Prestige.CurrentValue,
	}}

	err := s.Client.WriteJSON(session_message)
	if err != nil {
		fmt.Printf("failed to send message: %v\n", err)
		return
	}

	for {
		_, _, err := s.Client.NextReader()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				switch closeErr.Code {
				case websocket.CloseNormalClosure:
					fmt.Printf("client %d closed connection normally\n", s.Session.UserID)
				case websocket.CloseGoingAway:
					fmt.Printf("client %d is going away\n", s.Session.UserID)
				default:
					fmt.Printf("client %d closed with unexpected code %d: %v\n", s.Session.UserID, closeErr.Code, closeErr)
				}
			} else {
				fmt.Printf("client %d disconnected: %v\n", s.Session.UserID, err)
			}	
			break
		}
	}
}
