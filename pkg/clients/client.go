package clients

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type Client struct {
	Id       string
	Conn     net.Conn
	RecvChan chan string
	SendChan chan string
	DoneChan chan bool
	Wg       sync.WaitGroup
	RecvWg   sync.WaitGroup
	SendWg   sync.WaitGroup
	Session  *Session
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		Id:       conn.RemoteAddr().String(),
		Conn:     conn,
		RecvChan: make(chan string),
		SendChan: make(chan string),
		DoneChan: make(chan bool, 1),
		Session:  &Session{},
	}

	c.Wg.Add(1)
	c.RecvWg.Add(1)
	c.SendWg.Add(1)
	go c.receiveHandler(conn)
	go c.sendHandler()

	return c
}

func (c *Client) receiveHandler(conn net.Conn) {
	defer c.RecvWg.Done()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
				log.Println("Error reading from connection:", err)
			}
			c.DoneChan <- true
			return
		}

		data := string(buf[:n])
		log.Printf("[%s] <<< %s\n", c.Id, data)
		c.RecvChan <- data
	}
}

func (c *Client) sendHandler() {
	defer c.SendWg.Done()

	for {
		msg, ok := <-c.SendChan
		if !ok {
			return
		}

		log.Printf("[%s] >>> %s\n", c.Id, msg)
		if _, err := c.Conn.Write([]byte(msg)); err != nil {
			log.Println("Error writing to connection:", err)
			c.DoneChan <- true
			return
		}
	}
}

func (c *Client) Send(msg string) {
	c.SendChan <- msg
}

func (c *Client) Disconnect() {
	close(c.SendChan)
	c.SendWg.Wait()
	c.Conn.Close()
	c.RecvWg.Wait()
	close(c.RecvChan)
	close(c.DoneChan)
	c.Wg.Done()
	log.Println("Client disconnected:", c.Id)
}
