package commands

import (
	"fmt"
	"msnserver/pkg/clients"
)

func HandleSendILN(c *clients.Client, transactionID string, status string, email string, name string) {
	res := fmt.Sprintf("ILN %s %s %s %s\r\n", transactionID, status, email, name)
	c.Send(res)
}
