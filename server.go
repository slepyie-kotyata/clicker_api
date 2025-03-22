package main

import (
	"clicker_api/routes"

	"github.com/labstack/echo/v4"
)




func main() {
	e := echo.New()

	routes.InitEntryRoutes(e)

	e.Start(":1323")

}

