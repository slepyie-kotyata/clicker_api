package database

import (
	"clicker_api/environment"
	"clicker_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db_connection *gorm.DB

func GetDBConnection() *gorm.DB {
	if db_connection == nil {
		node_env := environment.GetVariable("NODE_ENV")
		if node_env == "development" {
			connectToSQLite()
		} else if node_env == "production" {
			connectToPostgres()
		} else {
			panic("undefined database connection")
		}
	}

	return db_connection
}

func connectDB (db *gorm.DB, err error) *gorm.DB {
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

func connectToSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	return connectDB(db, err)
}

func connectToPostgres() *gorm.DB {
	dsn := environment.GetVariable("DB_CONNECTION")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return connectDB(db, err)
}