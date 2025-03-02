package switchboard

import (
	"log"
	"msnserver/config"
	"msnserver/internal/switchboard"
	"msnserver/pkg/database"
	"msnserver/pkg/redis"
)

func main() {
	log.Println("Starting MSN switchboard server...")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	db, err := database.Load(c.Database)
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	rdb, err := redis.NewRedisClient(c.Redis)
	if err != nil {
		log.Fatalln("Error loading redis:", err)
	}

	ss := switchboard.NewSwitchboardServer(c, db, rdb)
	ss.Start()
}
