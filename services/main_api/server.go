package main

import (
	custommiddleware "clicker_api/services/main_api/custom_middleware"
	"clicker_api/services/main_api/database"
	"clicker_api/services/main_api/routes"
	"clicker_api/services/main_api/secret"
	"clicker_api/services/main_api/ws"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:4200",
			"https://clicker.enjine.ru",
			"https://enjine.ru",
    	},
    	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
    	AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
    	},
    	AllowCredentials: true,
	}))
	
	refresh := e.Group("/refresh")
	refresh.Use(custommiddleware.JWTMiddleware(secret.Refresh_secret))
	
	go ws.H.Run()
	go ws.P.Start()
	go database.A.Start()

	routes.InitEntryRoutes(e)
	routes.InitRefreshRoute(refresh)
	routes.InitWsRoutes(e)

	//gracefull shutdown
	go func(){
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			log.Printf("stopped listening: %v", err)
		}
	}()
		
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	<-ctx.Done()

	defer stop()
	
	log.Println("shutting down server..")

	database.A.Stop()
	ws.P.Stop()

	ctx_shd, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := e.Shutdown(ctx_shd); err != nil {
		log.Printf("shutting down with error: %v", err)
	} else {
		log.Println("shut down complete")
	}
}


