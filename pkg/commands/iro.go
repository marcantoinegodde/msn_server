package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendIRO(c *clients.Client, tid uint32, clients []*clients.Client) {
	for i, client := range clients {
		res := fmt.Sprintf("IRO %d %d %d %s %s\r\n", tid, i+1, len(clients), client.Session.Email, client.Session.DisplayName)
		c.Send(res)
	}
}
