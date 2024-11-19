package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/dispatch"
	"msnserver/pkg/database"
)

func main() {
	log.Println("Starting MSN dispatch server...")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	db, err := database.Load(c.Database)
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	ds := dispatch.NewDispatchServer(db, c)
	ds.Start()
}
