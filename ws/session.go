package ws

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/secret"
	"clicker_api/service"
	"clicker_api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

//TODO:разработать хранилище данных сессии

type SessionConn struct {
	session  		*models.Session
	client   		*websocket.Conn
	user_id    		uint
	messages 		chan Message
	done      		chan struct{}
}

func NewSession(conn *websocket.Conn) *SessionConn {
	return &SessionConn{
		session: 		nil,
		client: 		conn,
		user_id: 		0,
		messages: 		make(chan Message, 10),
		done: 			make(chan struct{}),
	}
}

const (
	write_wait = 10 * time.Second
	pong_wait = 60 * time.Second
	ping_period = (pong_wait * 9) / 10
	max_message_size = 10000
)

func (s *SessionConn) close() {
  fmt.Println("exiting session...")
	select {
	case <-s.done:
		return
	default:
		close(s.done)
		_ = s.client.Close()
	}
  fmt.Println("done!")
}

func (s *SessionConn) closeWithCode(code int, msg string) {
  _ = s.client.WriteControl(
      websocket.CloseMessage,
      websocket.FormatCloseMessage(code, msg),
      time.Now().Add(time.Second),
  )
  s.client.Close()
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
      var closeErr *websocket.CloseError
      
      if errors.As(err, &closeErr) {
        switch closeErr.Code {
        case websocket.CloseNormalClosure:
          log.Println("Normal close (1000)")
        case websocket.CloseGoingAway:
          log.Println("Going away (1001)")
        case websocket.CloseAbnormalClosure:
          log.Println("Abnormal close (1006)")
        default:
          log.Printf("Close code=%d text=%s", closeErr.Code, closeErr.Text)
        }
        return
      }

      log.Printf("Client disconnected: %v", err)
      return
    }

		fmt.Println("message has been recieved")

		var m Message
		if err = json.Unmarshal(message, &m); err != nil {
			log.Printf("invalid message error: %v", err)
			continue
		}

		switch m.MessageType {
		case Request:
			data, err := AuthorizeRequest(m.Data) 
			if err != nil{ 
				s.client.SetWriteDeadline(time.Now().Add(write_wait))
				message, _ := json.Marshal(map[string]interface{}{"message": err.Error()})

				hub.incoming <- HubEvent{
					Type: BroadcastToConnection,
					UserID: s.user_id,
					Session: s,
					Message: Message{
          				MessageType: Response, 
          				RequestID: m.RequestID,
          				RequestType: ErrorRequest,
          				Data: message,
        			},
				}
			}

			log.Println("message has been authorized")

			if s.user_id == 0 {
				s.user_id = utils.StringToUint(service.ExtractIDFromToken(data.Token, secret.Access_secret))
				hub.incoming <- HubEvent{
					Type:    RegisterConnection,
        			UserID:  s.user_id,
        			Session: s,
				}
			}
			s.InitAction(&m, data)
		default:
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
      		s.client.SetWriteDeadline(time.Now().Add(write_wait))

      		if !ok {
        	s.closeWithCode(websocket.CloseNormalClosure, "channel closed")
        		return
      		}

      		byte_message, err := utils.ToJSON(message)
      		if err != nil {
        		s.closeWithCode(websocket.CloseInternalServerErr, "encode error")
        		return
      		}

      		if err = s.client.WriteMessage(websocket.TextMessage, byte_message); err != nil {
        		return
      		}

		case <-ticker.C:
			s.client.SetWriteDeadline(time.Now().Add(write_wait))
			err := s.client.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				return
			}
			log.Println("ping sent")
		case <-s.done:
			return
		}

	}
}

func (s *SessionConn) InitAction(m *Message, data *RequestData) {
	switch m.RequestType {
	case SessionRequest:
    	s.session = database.InitSession(s.user_id)

    	data, _ := json.Marshal(map[string]interface{}{"session": s.session})

		hub.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: m.RequestType,
				Data:        data,
			},
		}

	default:
		return
	}
}

