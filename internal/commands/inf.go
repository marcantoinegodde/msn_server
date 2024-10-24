package commands

import (
	"fmt"
	"net"
)

func HandleINF(conn net.Conn, transactionID string) {
	conn.Write([]byte(fmt.Sprintf("INF %s MD5\r\n", transactionID)))
}
