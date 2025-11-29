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
		allowedOrigins := map[string]bool{
			"https://clicker.enjine.ru": 	true,
			"http://localhost:4200":     	true,
		}
	
		return allowedOrigins[origin]
	},
}

var session_manager = NewSessionManager()

func ServeWS(c echo.Context) error {
	fmt.Println("ORIGIN:", c.Request().Header.Get("Origin"))
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
