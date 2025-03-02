package dispatch

import (
	"log"
	"msnserver/config"
	"msnserver/internal/dispatch"
)

func main() {
	log.Println("Starting MSN dispatch server...")

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}

	ds := dispatch.NewDispatchServer(c)
	ds.Start()
}
