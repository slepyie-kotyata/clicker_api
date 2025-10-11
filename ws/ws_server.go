package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		
		allowedOrigins := map[string]bool{
			"wss://clicker.enjine.ru":    	true,
			"ws://localhost:4200":        	true,
			"https://clicker.enjine.ru": 	true,
			"http://localhost:4200":     	true,
		}

		return allowedOrigins[origin]
	},
}

