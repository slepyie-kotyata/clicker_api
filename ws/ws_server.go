package ws

import (
	"clicker_api/secret"
	"clicker_api/service"
	"clicker_api/utils"
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
			"wss://clicker.enjine.ru":    	true,
			"ws://localhost:4200":        	true,
			"https://clicker.enjine.ru": 	true,
			"http://localhost:4200":     	true,
		}

		return allowedOrigins[origin]
	},
}

func ServeWs(c echo.Context) error {
	token := c.QueryParam("token")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "6",
			"message": "missing token",
		})
	}
	
	err := service.ValidateToken(token, secret.Access_secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "6",
			"message": err.Error(),
		})
	}

	id := utils.StringToUint(service.ExtractIDFromToken(token, secret.Access_secret))

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	session_conn := NewSession(conn, id)
	
	go session_conn.readPump()
	go session_conn.writePump()

	return nil
}

