package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendJOI(c *clients.Client, clients []*clients.Client) {
	for _, client := range clients {
		res := fmt.Sprintf("JOI %s %s\r\n", c.Session.Email, c.Session.DisplayName)
		client.Send(res)
	}
}
