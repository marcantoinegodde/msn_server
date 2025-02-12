package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
)

func HandleSendILN(c *clients.Client, tid uint32, status database.Status, email string, name string) {
	res := fmt.Sprintf("ILN %d %s %s %s\r\n", tid, status, email, name)
	c.Send(res)
}
