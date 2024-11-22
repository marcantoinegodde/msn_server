package notification

import (
	"log"
	"net"
)

type Client struct {
	id       string
	conn     net.Conn
	sendChan chan string
}

func (c *Client) sendHandler() {
	defer c.conn.Close()

	for {
		msg, ok := <-c.sendChan
		if !ok {
			log.Printf("Channel closed: %s\n", c.id)
			return
		}

		log.Println(">>>", msg)
		if _, err := c.conn.Write([]byte(msg)); err != nil {
			log.Println("Error:", err)
			return
		}
	}
}
