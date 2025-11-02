package ws

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/utils"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type SessionConn struct {
	session  		*models.Session
	client   		*websocket.Conn
	messages 		chan Message
	done      		chan struct{}
}

func NewSession(conn *websocket.Conn, id uint) *SessionConn {
	return &SessionConn{
		session: 	database.InitSession(id),
		client: 	conn,
		messages: 	make(chan Message),
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

	for {
		_, message, err := s.client.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			return
		}

		s.client.SetReadDeadline(time.Now().Add(pong_wait))

		var m Message
		err = json.Unmarshal(message, &m.Data)
		if err != nil {
			log.Printf("invalid message error: %v", err)
			continue
		}

		s.messages <- m
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
			_ = s.client.SetWriteDeadline(time.Now().Add(write_wait))
			if !ok {
				_ = s.client.WriteMessage(websocket.CloseMessage, []byte{})
			}

			byte_message, err := utils.ToJSON(message)
			if err != nil {
				return
			}

			err = s.client.WriteMessage(websocket.TextMessage, byte_message)
			if err != nil {
				return
			}

		case <-ticker.C:
			s.client.SetWriteDeadline(time.Now().Add(write_wait))
			data, _ := json.Marshal(map[string]string{"data": "keep alive"})
			message, _ := utils.ToJSON(Message{MessageType: KeepAlive, Data: data})
			
			err := s.client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

			s.client.SetReadDeadline(time.Now().Add(pong_wait))

		case <-s.done:
			return
		}

	}
}

