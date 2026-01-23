package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		if origin == "origin" || origin == ""{
			return true
		}
		
		allowedOrigins := map[string]bool{
			"https://clicker.enjine.ru": 	true,
			"http://localhost:4200":     	true,
		}

		return allowedOrigins[origin]
	},
}

func ServeWs(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	log.Println("connection upgraded")

	session_conn := NewSession(conn)


	log.Println("session created")
	
	go session_conn.readPump()
	go session_conn.writePump()

	return nil
}

