package ws

import (
	"clicker_api/secret"
	"clicker_api/service"
	"encoding/json"
	"errors"
)

type MessageType string
const (
	Response = "response"
	Request = "request"
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

type RequestData struct {
	Token  	string 		`json:"token"`
	Action 	RequestType	`json:"request_type"`
	Param 	int			`json:"param,omitempty"`
}

func AuthorizeRequest(request_data json.RawMessage) (*RequestData, error) {
	var data RequestData
	if err := json.Unmarshal(request_data, &data); err != nil {
		return nil, errors.New("invalid data")
	}

	if data.Token == "" {
 		return nil, errors.New("missing token")
	}

	if err := service.ValidateToken(data.Token, secret.Access_secret); err != nil {
		return nil, err
	}

	return &data, nil
}