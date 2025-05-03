package handlers

import (
	"clicker_api/database"
	"clicker_api/environment"
)

var (
	DB = database.GetDBConnection()
	Secret = environment.GetVariable("ACCESS_TOKEN_SECRET")
)