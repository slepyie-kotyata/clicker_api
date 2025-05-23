package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"fmt"
	"math"
	"net"
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
	Success  chan struct{}
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
	}()

	wg.Wait()

	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").First(&s.Session, s.Session.ID)

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
			s.Client.SetReadDeadline(time.Now().Add(time.Duration(seconds_interval) * time.Second))

			select {
			case <- s.Success:
				s.Client.SetReadDeadline(time.Time{})
			case <- time.After(time.Duration(seconds_interval) * time.Second):
				fmt.Printf("client %d did not reply in time\n", s.Session.UserID)
				s.Client.Close()
				return
			case <- s.Done:
				return
			}
		case <- s.Done:
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
				err := s.Client.WriteJSON(session_message)
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

	go func() {
		for {
			_, message, err := s.Client.ReadMessage()
			if err != nil {
				close(done)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("client %d disconnected unexpectedly: %v\n", s.Session.UserID, err)
				} else if err, ok := err.(net.Error); ok && err.Timeout() {
					fmt.Printf("client %d did not respond in time (timeout)\n", s.Session.UserID)
				} else {
					fmt.Printf("client %d read error: %v\n", s.Session.UserID, err)
				}
				return
			}

			if string(message) == `"success"` {
				select {
				case s.Success <- struct{}{}:
				default:
				}
			}
		}
	}()
	<-done
}
