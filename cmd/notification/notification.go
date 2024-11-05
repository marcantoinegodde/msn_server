package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/notification"
	"msnserver/pkg/database"
)

func main() {
	log.Println("Starting MSN notification server...")

	config.LoadConfig()

	db, err := database.Load()
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	notification.StartNotificationServer(db)
}
