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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Europe/Paris", c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")

	return db, nil
}

func ResetUsersStatus(db *gorm.DB) {
	// TODO: remove this function, store users' status in cache instead

	err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("status", FLN).Error
	if err != nil {
		log.Fatalf("Failed to reset column 'status': %v", err)
	}

	log.Println("Column 'status' reseted successfully")
}
