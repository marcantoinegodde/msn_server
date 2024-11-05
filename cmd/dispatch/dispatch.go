package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/dispatch"
	"msnserver/pkg/database"
)

func main() {
	log.Println("Starting MSN dispatch server...")

	config.LoadConfig()

	db, err := database.Load()
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	dispatch.StartDispatchServer(db)
}
