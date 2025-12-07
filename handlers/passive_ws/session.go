package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"fmt"
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
	Session  		models.Session
	Client   		*websocket.Conn
	Messages 		chan SessionMessage
	Done     		chan struct{}
	Success  		chan struct{}
	Closed   		bool
	lastMessage 	SessionMessage
	last_mu      	sync.RWMutex
	mu       		sync.RWMutex
}

var seconds_interval uint = 3

func (s *Session) createMessage() {
	var fresh models.Session
	database.DB.Preload("Level").Preload("Prestige").Where("id = ?", s.Session.ID).First(&fresh)

	s.last_mu.Lock()
	s.lastMessage = SessionMessage{
		Money:           fresh.Money,
		Dishes:          fresh.Dishes,
		Rank:            fresh.Level.Rank,
		XP:              fresh.Level.XP,
		PrestigeCurrent: fresh.Prestige.CurrentValue,
	}
	s.last_mu.Unlock()
}

func (s *Session) UpdateSessionState(seconds uint) {
	upgrade_stats := service.CountBoostValues(service.FilterUpgrades(&s.Session, true))

	if upgrade_stats.MpS == 0 && upgrade_stats.DpS == 0 {
		s.createMessage()
		return
	}

	prestige_boost := s.Session.Prestige.AccumulatedValue

	if prestige_boost == 0 {
		prestige_boost = 1
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		s.PrestigeUpgrade(upgrade_stats, seconds)
		
	}()

	go func() {
		defer wg.Done()
		s.PassiveSellUpdate(upgrade_stats, seconds, prestige_boost)
	}()

	go func() {
		defer wg.Done()
		s.PassiveCookUpdate(upgrade_stats, seconds, prestige_boost)
	}()

	wg.Wait()

	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").First(&s.Session, s.Session.ID)

	if s.Closed {
		return
	}

	s.createMessage()
}

func (s *Session) StartPassiveLoop() {
	update_ticker := time.NewTicker(time.Duration(seconds_interval) * time.Second)
	send_ticker := time.NewTicker(time.Duration(seconds_interval) * time.Second)

	defer func() {
		update_ticker.Stop()
		send_ticker.Stop()
	}()

	const grace_period = 2 * time.Second

	go func() {
		for {
			select {
			case <- update_ticker.C:
				s.UpdateSessionState(seconds_interval)
				s.Client.SetReadDeadline(time.Now().Add(time.Duration(seconds_interval) * time.Second + grace_period))

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
	}()

	for {
		select {
		case <-send_ticker.C:
			s.createMessage()
			s.last_mu.RLock()
			msg := s.lastMessage
			s.last_mu.RUnlock()
			s.Messages <- msg
		case <-s.Done:
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
