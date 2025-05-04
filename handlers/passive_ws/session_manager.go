package passivews

import (
	"clicker_api/handlers"
	"clicker_api/models"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type SessionManager struct {
	Sessions map[uint]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[uint]*Session),
	}
}

func (sm *SessionManager) CloseSession(id uint) {
	sm.mu.RLock()
	session, ok := sm.Sessions[id]
	sm.mu.RUnlock()
	if !ok {
		return
	}

	session.Client.Close()
	close(session.Messages)
	sm.mu.Lock()
	delete(sm.Sessions, id)
	sm.mu.Unlock()
}

var count int = 0

func (sm *SessionManager) CreateAndAddToSession(conn *websocket.Conn, id uint) error {
	var this_session models.Session
	handlers.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&this_session)

	if this_session.ID == 0 {
		return errors.New("game is not initialized")
	}

	session := Session{
		Session: this_session,
		Client: conn,
		Messages: make(chan SessionMessage),
	}

	sm.Sessions[this_session.ID] = &session

	go session.HandleConnection(sm)
	fmt.Printf("goroutine number %d started\n", count)
	count++
	return nil
}