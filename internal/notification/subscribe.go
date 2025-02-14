package notification

import (
	"context"
	"encoding/json"
	"log"
	"msnserver/pkg/commands"
)

func (ns *NotificationServer) subscribe() {
	pubsub := ns.rdb.Subscribe(context.TODO(), ns.config.Redis.PubSubChannel)

	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var rngMsg commands.RNGMessage

		if err := json.Unmarshal([]byte(msg.Payload), &rngMsg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		log.Printf("Received message from switchboard: %+v\n", rngMsg)

		if err := commands.HandleRNG(ns.rdb, ns.m, ns.clients, rngMsg); err != nil {
			log.Println("Error handling RNG message:", err)
		}
	}
}
