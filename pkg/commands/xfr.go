package commands

import (
	"fmt"
	"msnserver/config"
)

func HandleXFR(c chan string, cf config.DispatchServer, transactionID string) {
	res := fmt.Sprintf("XFR %s NS %s:%s 0 %s:%s\r\n", transactionID, cf.NotificationServerAddr, cf.NotificationServerPort, cf.ServerAddr, cf.ServerPort)
	c <- res
}
