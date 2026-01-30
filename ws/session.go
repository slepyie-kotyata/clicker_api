package ws

import (
	"clicker_api/database"
	"clicker_api/secret"
	"clicker_api/service"
	"clicker_api/utils"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type SessionConn struct {
	client   		*websocket.Conn
	user_id    		uint
	messages 		chan Message
	done      		chan struct{}
}

func NewSession(conn *websocket.Conn) *SessionConn {
	return &SessionConn{
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
  	log.Println("exiting session...")
	if s.user_id != 0 {
		log.Println("unregistered from hub")
		
		H.incoming <- HubEvent{
			Type:    UnregisterConnection,
        	UserID:  s.user_id,
        	Session: s,
    	}
	}

	select {
	case <-s.done:
		return
	default:
		close(s.done)
		_ = s.client.Close()
	}
  	log.Println("done!")
}

func (s *SessionConn) writeCloseMessage(code int, msg string) {
  	_ = s.client.WriteControl(
    	websocket.CloseMessage,
      	websocket.FormatCloseMessage(code, msg),
      	time.Now().Add(time.Second),
  	)
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

		log.Println("message has been recieved")

		var m Message
		if err = json.Unmarshal(message, &m); err != nil {
			log.Printf("invalid message error: %v", err)
			continue
		}

		switch m.MessageType {
		case Request:
			data, err := AuthorizeRequest(m.Data) 
			if err != nil { 
				s.client.SetWriteDeadline(time.Now().Add(write_wait))
				message, _ := json.Marshal(map[string]interface{}{"message": err.Error()})

				byte_message, _ := json.Marshal(Message{
					MessageType: Response, 
					RequestID: m.RequestID,
					RequestType: ErrorRequest,
					Data: message,
				})

				if err = s.client.WriteMessage(websocket.TextMessage, byte_message); err != nil {
					return
				}
				continue
			}

			if s.user_id == 0 {
				s.user_id = utils.StringToUint(service.ExtractIDFromToken(data.Token, secret.Access_secret))
				H.incoming <- HubEvent{
					Type:    RegisterConnection,
        			UserID:  s.user_id,
					Session: s,
				}
			}
			log.Println("authorized")
			log.Println("listening for actions")
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
        	s.writeCloseMessage(websocket.CloseNormalClosure, "channel closed")
        		return
      		}

      		byte_message, err := json.Marshal(message)
      		if err != nil {
        		s.writeCloseMessage(websocket.CloseInternalServerErr, "encode error")
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
		s_bd := database.InitSession(s.user_id)
		session := database.CreateSessionState(s_bd)
		
    	data, _ := json.Marshal(map[string]interface{}{"session": SessionResponse{
			UserID: s_bd.UserID,
			UserEmail: s_bd.UserEmail,
			Money: s_bd.Money,
			Dishes: s_bd.Dishes,
			Level: struct {
				Rank uint    `json:"rank"`
				XP   float64 `json:"xp"`
			}{
				Rank: s_bd.Level.Rank,
				XP:   s_bd.Level.XP,
			},
			Prestige: struct {
				CurrentValue       float64 `json:"current_value"`
				CurrentBoostValue  float64 `json:"current_boost_value"`
				AccumulatedValue   float64 `json:"accumulated_value"`
			}{
				CurrentValue:      s_bd.Prestige.CurrentValue,
				CurrentBoostValue: s_bd.Prestige.CurrentBoostValue,
				AccumulatedValue:  s_bd.Prestige.AccumulatedValue,
			},
			Upgrades: struct {
				Available []service.FilteredUpgrade `json:"available"`
				Current   []service.FilteredUpgrade `json:"current"`
			}{
				Available: service.FilterUpgrades(session, false),
				Current: service.FilterUpgrades(session, true),
			},
		}})

		H.incoming <- HubEvent{
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
	case BuyRequest:
		response, r_type := s.Buy(data.Param)

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case CookRequest:
		response, r_type := s.Cook()

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case SellRequest:
		response, r_type := s.Sell()

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case ListRequest:
		response, r_type := s.ListUpgrades()

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case LevelUpRequest:
		response, r_type := s.LevelUp()

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case CheckLevelRequest:
		response, r_type := s.GetLevel()

		data, _ := json.Marshal(response)
		
		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	case ResetRequest:
		response, r_type := s.ResetSession()

		data, _ := json.Marshal(response)

		H.incoming <- HubEvent{
			Type: BroadcastToConnection,
			UserID: s.user_id,
			Session: s,
			Message: Message{
				MessageType: Response,
				RequestID:   m.RequestID,
				RequestType: r_type,
				Data:        data,
			},
		}
	default:
		return
	}
}

