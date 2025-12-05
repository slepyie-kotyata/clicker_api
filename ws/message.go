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
	ErrorRequest = "error"
)

type Message struct {
	MessageType	MessageType		`json:"message_type"`
	RequestID	string			`json:"request_id"`
	RequestType RequestType		`json:"request_type"`
	Data 		json.RawMessage	`json:"data"`
}

type RequestData struct {
	Token  		string 		`json:"token"`
	Param 		int			`json:"param,omitempty"`
}

type SessionResponse struct {
	UserID     		uint		`json:"user_id"`
	UserEmail 		string 		`json:"user_email"`
	Money      		uint		`json:"money"`
	Dishes     		uint		`json:"dishes"`
	Level      		struct {
		Rank		uint	`json:"rank"`	
		XP			float64	`json:"xp"`
	}	`json:"level"`
	Prestige   		struct {
		CurrentValue 		float64		`json:"current_value"`
		CurrentBoostValue 	float64		`json:"current_boost_value"`
		AccumulatedValue   	float64		`json:"accumulated_value"`
	}   `json:"prestige"`
	Upgrades   		struct {
		Available	[]service.FilteredUpgrade	`json:"available"`
		Current 	[]service.FilteredUpgrade	`json:"current"`
	}   `json:"upgrades"`
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