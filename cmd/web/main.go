package web

import (
	"log"
	"msnserver/config"
	"msnserver/internal/web"
	"msnserver/pkg/database"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	db, err := database.Load(c.Database)
	if err != nil {
		log.Fatalln("Error loading database:", err)
	}

	ws := web.NewWebServer(c, db)
	ws.Start()
}
