package main

import (
	"log"
	"msnserver/config"
	"msnserver/internal/switchboard"
)

func main() {
	log.Println("Starting MSN switchboard server...")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	ss := switchboard.NewSwitchboardServer(c)
	ss.Start()
}
