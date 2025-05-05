package handlers

import (
	"clicker_api/environment"
)

var (
	Secret = environment.GetVariable("ACCESS_TOKEN_SECRET")
)