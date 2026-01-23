package database

import (
	"clicker_api/pkg/models"
	"clicker_api/services/main_api/environment"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db_connection *gorm.DB

func GetDBConnection() *gorm.DB {
	if db_connection == nil {
		connectDB()
	}
	return db_connection
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(environment.GetVariable("DB_CONNECTION")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(
		&models.User{}, 
		&models.Password{},
		&models.Session{},
		&models.Upgrade{},
		&models.Boost{},
		&models.SessionUpgrade{},
		&models.LevelXP{},
		&models.Level{},
		&models.Prestige{},
	)

	_ = db.SetupJoinTable(&models.Session{}, "Upgrades", &models.SessionUpgrade{})

	db_connection = db

	return db_connection
}

func ConnectRedis() *redis.Client {
	opt, err := redis.ParseURL(environment.GetVariable("REDIS_CONNECTION"))
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}