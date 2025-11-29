package passivews

import (
	"clicker_api/secret"
	"clicker_api/service"
	"clicker_api/utils"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		fmt.Println("ORIGIN:", origin)
		allowedOrigins := map[string]bool{
			"https://clicker.enjine.ru": 	true,
			"http://localhost:4200":     	true,
		}
	
		return allowedOrigins[origin]
	},
}

var session_manager = NewSessionManager()

func ServeWS(c echo.Context) error {
	token := c.Request().Header.Get("Sec-Websocket-Protocol")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "6",
			"message": "missing token",
		})
	}
	fmt.Println("TOKEN RECEIVED:", token)

	err := service.ValidateToken(token, secret.Access_secret)
	if err != nil {
		fmt.Println("ValidateToken error:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "6",
			"message": err.Error(),
		})
	}
	
	id := utils.StringToUint(service.ExtractIDFromToken(token, secret.Access_secret))
	fmt.Println("EXTRACTED USER ID:", id)

	fmt.Println("UPGRADING WS...")
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("UPGRADE RESULT:", err)
		c.Logger().Error(err)
		return err
	}

	err = session_manager.CreateAndAddToSession(conn, id)

	if err != nil {
		error_message := map[string]interface{}{
			"message": err.Error(),
		}
		_ = conn.WriteJSON(error_message)
		conn.Close()
		return nil
	}

	return nil
}
