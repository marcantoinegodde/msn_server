package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendACK(c *clients.Client, tid uint32) {
	res := fmt.Sprintf("ACK %d\r\n", tid)
	c.Send(res)
}
