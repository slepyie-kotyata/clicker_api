package handlers

import (
	"clicker_api/environment"
)

var (
	Access_secret = environment.GetVariable("ACCESS_TOKEN_SECRET")
	Refresh_secret = environment.GetVariable("REFRESH_TOKEN_SECRET")
)