package database

var (
	DB = GetDBConnection()
	RClient = ConnectRedis()
)