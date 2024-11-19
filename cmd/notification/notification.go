package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/notification"
	"msnserver/pkg/database"
)

func main() {
	log.Println("Starting MSN notification server...")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	db, err := database.Load(c.Database)
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	ns := notification.NewNotificationServer(db, c)
	ns.Start()
}
