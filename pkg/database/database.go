package database

import (
	"fmt"
	"log"
	"msnserver/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Load() (*gorm.DB, error) {
	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Paris", config.Config.Database.Host, config.Config.Database.User, config.Config.Database.Password, config.Config.Database.DBName, config.Config.Database.Port, config.Config.Database.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return db, nil
}
