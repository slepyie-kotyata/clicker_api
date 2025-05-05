package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"fmt"
	"math"
	"sync"
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
	Done     chan struct{}
	Closed   bool
	mu       sync.RWMutex
}

var seconds_interval uint = 3

func (s *Session) UpdateSessionState(seconds uint) {
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(s.Session, true))
	current_prestige := math.Round(1 + 0.05 * s.Session.PrestigeValue)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		s.PrestigeUpgrade(upgrade_stats, seconds)
		
	}()

	go func() {
		defer wg.Done()
		s.PassiveSellUpdate(upgrade_stats, seconds, current_prestige)
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Before Update:", s.Session.Dishes)
		s.PassiveCookUpdate(upgrade_stats, seconds, current_prestige)
	}()

	wg.Wait()

	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").First(&s.Session, s.Session.ID)
	fmt.Println("After Upgrade:", s.Session.Dishes)

	if s.Closed {
		return
	}

	s.Messages <- SessionMessage{
		Money: s.Session.Money,
		Dishes: s.Session.Dishes,
		Rank: s.Session.Level.Rank,
		XP: s.Session.Level.XP,
		PrestigeCurrent: s.Session.Prestige.CurrentValue,
	}
}

func (s *Session) StartPassiveLoop() {
	ticker := time.NewTicker(time.Duration(seconds_interval) * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <- ticker.C:
			s.UpdateSessionState(seconds_interval)
		case <- s.Done:
			fmt.Println("stopped")
			return
		}
	}
}

func (s *Session) HandleConnection(sm *SessionManager) {
	defer sm.CloseSession(s.Session.ID)

	done := make(chan struct{})

	go func() {
		for {
			select {
			case session_message := <- s.Messages:
				err := s.Client.WriteJSON(map[string]interface{}{"message": session_message})
				if err != nil {
					fmt.Println("failed to send message: %v\n", err)
					close(done)
					return
				}
			case <- done:
				return
			}
		}
	}()

	for {
		_, _, err := s.Client.NextReader()
		if err != nil {
			close(done)
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
