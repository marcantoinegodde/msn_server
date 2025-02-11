package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendILN(c *clients.Client, tid uint32, status string, email string, name string) {
	res := fmt.Sprintf("ILN %d %s %s %s\r\n", tid, status, email, name)
	c.Send(res)
}
