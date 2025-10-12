package ws

import "encoding/json"

const CookAction = "cook"
const SellAction = "sell"
const BuyAction = "buy"
const ListAction = "list"
const SessionAction = "session"
const LevelUpAction = "level_up"
const CheckLevelAction = "check_level"
const ResetAction = "reset"
const PassiveAction = "passive"
const LeaveAction = "leave"

type Message struct {
	Action 	string
	Data 	json.RawMessage
}