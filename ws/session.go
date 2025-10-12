package ws

import (
	"clicker_api/database"
	"clicker_api/models"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)


func InitSession(id uint) *models.Session {
	var (
		session models.Session
		user models.User
	)
	
	database.DB.Preload("Prestige").Preload("Level").Where("user_id = ?", id).First(&session)
	database.DB.Select("email").First(&user, id)
	
	if session.ID > 0 {
		return &session
	}
	
	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		PrestigeValue: 0,
		PrestigeBoost: 0,
		UserID: id,
		UserEmail: user.Email,
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
	
	database.DB.Preload("Prestige").Preload("Level").Where("user_id = ?", id).First(&new_session)
	
	return &new_session
}

type SessionConn struct {
	session  		*models.Session
	client   		*websocket.Conn
	messages 		chan *Message
	done      		chan struct{}
}

func NewSession(conn *websocket.Conn, id uint) *SessionConn {
	return &SessionConn{
		session: 	InitSession(id),
		client: 	conn,
		messages: 	make(chan *Message),
		done: 		make(chan struct{}),

	}
}

const (
	write_wait = 10 * time.Second

	pong_wait = 60 * time.Second

	ping_period = (pong_wait * 9) / 10

	max_message_size = 10000
)

func (s *SessionConn) close() {
	select {
	case <-s.done:
		return
	default:
		close(s.done)
		_ = s.client.Close()
	}
}

func (s *SessionConn) readPump() {
	defer s.close()

	s.client.SetReadLimit(max_message_size)
	s.client.SetReadDeadline(time.Now().Add(pong_wait))
	s.client.SetPongHandler(func(string) error{
		log.Println("✅ Pong received from client")
		s.client.SetReadDeadline(time.Now().Add(pong_wait))
		return nil
	})

	for {
		_, message, err := s.client.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			return
		}

		var m Message
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Printf("invalid message error: %v", err)
			continue
		}

	}
}

func (s *SessionConn) writePump() {
	ticker := time.NewTicker(ping_period)
	defer func() {
		ticker.Stop()
		s.client.Close()
	}()

	for {
		select {
		case message, ok := <-s.messages:
			_ = s.client.SetReadDeadline(time.Now().Add(write_wait))
			if !ok {
				_ = s.client.WriteMessage(websocket.CloseMessage, []byte{})
			}

			err := s.client.WriteMessage(websocket.TextMessage, message.Data)
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = s.client.SetWriteDeadline(time.Now().Add(write_wait))
			err := s.client.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		case <-s.done:
			return
		}

	}
}

