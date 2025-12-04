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

var hub = NewHub()

func ServeWs(c echo.Context) error {
	// token := c.Request().Header.Get("Sec-Websocket-Protocol")

	// if token == "" {
	// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
	// 		"status": "6",
	// 		"message": "missing token",
	// 	})
	// }
	
	// err := service.ValidateToken(token, secret.Access_secret)
	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
	// 		"status": "6",
	// 		"message": err.Error(),
	// 	})
	// }

	// id := utils.StringToUint(service.ExtractIDFromToken(token, secret.Access_secret))

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	log.Println("connection upgraded")

	session_conn := NewSession(conn)
	// data := map[string]interface{}{
	// 	"request_type": SessionRequest,
	// 	"session": session_conn.session,
	// }

	log.Println("session created")

	// m_data, err := utils.ToJSON(data)
	// if err != nil {
	// 	log.Println("failed to initialize session", err)
	// 	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "init failed"))
	// 	conn.Close()
	// 	return nil
	// }

	// log.Println("session has been initialized!")
	
	go session_conn.readPump()
	go session_conn.writePump()

	go hub.Run()

	// session_conn.messages <- Message{MessageType: Response, Data: m_data}

	return nil
}

