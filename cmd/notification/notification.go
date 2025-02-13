package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/notification"
	"msnserver/pkg/database"
	"msnserver/pkg/redis"
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
	database.ResetUsersStatus(db)

	rdb, err := redis.NewRedisClient(c.Redis)
	if err != nil {
		log.Fatalln("Error loading redis:", err)
	}

	ns := notification.NewNotificationServer(c, db, rdb)
	ns.Start()
}
