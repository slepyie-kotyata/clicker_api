package secret

import "clicker_api/services/main_api/environment"

var (
	Access_secret = environment.GetVariable("ACCESS_TOKEN_SECRET")
	Refresh_secret = environment.GetVariable("REFRESH_TOKEN_SECRET")
)