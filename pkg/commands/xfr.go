package commands

import (
	"fmt"
	"msnserver/config"
	"msnserver/pkg/clients"
)

func HandleXFR(cf config.DispatchServer, c *clients.Client, transactionID string) {
	res := fmt.Sprintf("XFR %s NS %s:%s\r\n", transactionID, cf.NotificationServerAddr, cf.NotificationServerPort)
	c.SendChan <- res
}
