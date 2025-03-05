package sessions

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendBYE(c *clients.Client, clients []*clients.Client) {
	for _, client := range clients {
		res := fmt.Sprintf("BYE %s\r\n", c.Session.Email)
		client.Send(res)
	}
}
