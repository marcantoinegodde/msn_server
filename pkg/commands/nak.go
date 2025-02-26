package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendNAK(c *clients.Client, tid uint32) {
	res := fmt.Sprintf("NAK %d\r\n", tid)
	c.Send(res)
}
