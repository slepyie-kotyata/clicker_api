package ws

import "encoding/json"

type ActionType string
const (
	CookAction = "action_cook"
	SellAction = "action_sell"
	BuyAction = "upgrade_buy"
	ListAction = "upgrade_list"
	SessionAction = "session"
	LevelUpAction = "level_up"
	CheckLevelAction = "level_check"
	ResetAction = "session_reset"
	PassiveAction = "action_passive"
	LeaveAction = "action_leave"
)

type Message struct {
	Action 	ActionType		`json:"action"`
	Data 	json.RawMessage	`json:"data"`
}

func NewMessage(data json.RawMessage, action ActionType) Message {
	return Message{
		Action: action, 
		Data: data,
	}
}