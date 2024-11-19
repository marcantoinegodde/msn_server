package commands

import (
	"fmt"
	"log"
	"msnserver/config"
	"net"
)

func HandleXFR(conn net.Conn, c config.DispatchServer, transactionID string) {
	res := fmt.Sprintf("XFR %s NS %s:%s 0 %s:%s\r\n", transactionID, c.NotificationServerAddr, c.NotificationServerPort, c.ServerAddr, c.ServerPort)
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
