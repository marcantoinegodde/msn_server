package clients

import (
	"log"
	"net"
	"sync"
)

type Client struct {
	Id       string
	Conn     net.Conn
	SendChan chan string
	Wg       sync.WaitGroup
	Session  *Session
}

func (c *Client) SendHandler() {
	defer c.Conn.Close()

	for {
		msg, ok := <-c.SendChan
		if !ok {
			log.Printf("Channel closed: %s\n", c.Id)
			return
		}

		log.Printf("[%s] >>> %s\n", c.Id, msg)
		if _, err := c.Conn.Write([]byte(msg)); err != nil {
			log.Println("Error:", err)
			return
		}
	}
}
