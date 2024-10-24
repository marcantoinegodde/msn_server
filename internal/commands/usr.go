package commands

import (
	"fmt"
	"msnserver/config"
	"net"
)

func HandleUSRDispatch(conn net.Conn, transactionID string) {
	conn.Write([]byte(fmt.Sprintf("XFR %s NS %s:%s 0 %s:%s\r\n", transactionID, config.Config.DispatchServer.NotificationServerAddr, config.Config.DispatchServer.NotificationServerPort, config.Config.DispatchServer.ServerAddr, config.Config.DispatchServer.ServerPort)))
}
