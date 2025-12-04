package ws

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

func NewHub() *Hub {
    return &Hub{
        incoming: make(chan HubEvent, 256),
        sessions: make(map[uint]map[*SessionConn]bool),
    }
}

func (h *Hub) Run() {
    for event := range h.incoming {
        switch event.Type {
        case RegisterConnection:
            if h.sessions[event.UserID] == nil {
                h.sessions[event.UserID] = make(map[*SessionConn]bool)
            }
            h.sessions[event.UserID][event.Session] = true

        case BroadcastToConnection:
            if sessions, ok := h.sessions[event.UserID]; ok {
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
                }
            }
        }
    }
}

