package ws

import "encoding/json"

type MessageType string
const (
	Response = "response"
	Request = "request"
	KeepAlive = "keep_alive"
)

type RequestType string
const (
	CookRequest = "cook"
	SellRequest = "sell"
	BuyRequest = "upgrade_buy"
	ListRequest = "upgrade_list"
	SessionRequest = "session"
	LevelUpRequest = "level_up"
	CheckLevelRequest = "level_check"
	ResetRequest = "session_reset"
	PassiveRequest = "passive"
	LeaveRequest = "leave"
)

type Message struct {
	MessageType	MessageType		`json:"message_type"`
	Data 		json.RawMessage	`json:"data"`
}