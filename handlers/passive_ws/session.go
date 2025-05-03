package passivews

import (
	"clicker_api/models"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type SessionMessage struct {
	Money  uint 	`json:"money"`
	Dishes uint 	`json:"dishes"`
	Rank   uint 	`json:"rank"`
	XP     float64  `json:"xp"`
}

type Session struct {
	Session  models.Session
	Client   *websocket.Conn
	Messages chan SessionMessage
	once     sync.Once
}

func (s *Session) HandleConnection(sm *SessionManager) {
	defer sm.CloseSession(s.Session.ID)

	session_message := map[string]interface{}{"message":SessionMessage{
		Money: s.Session.Money,
		Dishes: s.Session.Dishes,
		Rank: s.Session.Level.Rank,
		XP: s.Session.Level.XP,
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
