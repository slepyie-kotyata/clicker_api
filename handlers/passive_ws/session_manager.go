package passivews

import (
	"clicker_api/database"
	"clicker_api/models"
	"errors"
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

	if session.Client != nil {
		_ = session.Client.Close()
	}

	session.mu.Lock()
	if !session.Closed {
		session.Closed = true
		close(session.Done)
		close(session.Messages)
	}
	session.mu.Unlock()

	sm.mu.Lock()
	delete(sm.Sessions, id)
	sm.mu.Unlock()
}

func (sm *SessionManager) CreateAndAddToSession(conn *websocket.Conn, id uint) error {
	var this_session models.Session
	database.DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&this_session)

	if this_session.ID == 0 {
		return errors.New("game is not initialized")
	}

	_, ok := sm.Sessions[this_session.ID]

	if ok {
		return errors.New("session is already running")
	}

	session := Session{
		Session: this_session,
		Client: conn,
		Messages: make(chan SessionMessage),
		Done: make(chan struct{}),
		Success: make(chan struct{}, 1),
	}

	sm.Sessions[this_session.ID] = &session

	go session.HandleConnection(sm)
	go session.StartPassiveLoop()

	return nil
}