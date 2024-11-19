package database

import (
	"fmt"
	"log"
	"msnserver/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Load(c config.Database) (*gorm.DB, error) {
	log.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Paris", c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return db, nil
}
