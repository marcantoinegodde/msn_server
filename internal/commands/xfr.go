package commands

import (
	"fmt"
	"log"
	"msnserver/config"
	"net"
)

func HandleXFR(conn net.Conn, transactionID string) {
	res := fmt.Sprintf("XFR %s NS %s:%s 0 %s:%s\r\n", transactionID, config.Config.DispatchServer.NotificationServerAddr, config.Config.DispatchServer.NotificationServerPort, config.Config.DispatchServer.ServerAddr, config.Config.DispatchServer.ServerPort)
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
