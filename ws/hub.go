package ws

import (
	"clicker_api/database"
	"fmt"
	"log"
)

type Hub struct {
    incoming chan HubEvent
    sessions map[uint]map[*SessionConn]bool
}

type HubEventType string
const (
    RegisterConnection = "register"
    UnregisterConnection = "unregister"
    BroadcastToConnection = "broadcast"
)

type HubEvent struct {
    Type      HubEventType
    UserID    uint
    Session   *SessionConn
    Message   Message
}

var H = NewHub()

func NewHub() *Hub {
    return &Hub{
        incoming: make(chan HubEvent, 256),
        sessions: make(map[uint]map[*SessionConn]bool),
    }
}

func (h *Hub) GetActiveUsers() []uint {
    users := make([]uint, 0, len(h.sessions))
    for id := range h.sessions {
        users = append(users, id)
    }
    return users
}

func (h *Hub) Run() {
    for event := range h.incoming {
        switch event.Type {
        case RegisterConnection:
            if h.sessions[event.UserID] == nil {
                h.sessions[event.UserID] = make(map[*SessionConn]bool)
            }
            h.sessions[event.UserID][event.Session] = true

            fmt.Println(len(h.sessions[event.UserID]))

        case BroadcastToConnection:
            if sessions, ok := h.sessions[event.UserID]; ok {
                log.Println("got the message!")
                for s := range sessions {
                    select {
                    case s.messages <- event.Message:
                    default:
                        close(s.messages)
                        delete(sessions, s)
                    }
                }
            }

        case UnregisterConnection:
            if sessions, ok := h.sessions[event.UserID]; ok {
                delete(sessions, event.Session)
                if len(sessions) == 0 {
                    delete(h.sessions, event.UserID)
                    database.SetTTL(event.UserID)
                }
            }
            fmt.Println(len(h.sessions[event.UserID]))
        }
    }
}

