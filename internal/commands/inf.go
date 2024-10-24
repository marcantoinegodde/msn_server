package commands

import (
	"fmt"
	"net"
	"strings"
)

var supportedAuthMethods = []string{"MD5"}

func HandleINF(conn net.Conn, arguments string) {
	transactionID, _, err := parseTransactionID(arguments)
	if err != nil {
		conn.Close()
		return
	}

	conn.Write([]byte(fmt.Sprintf("INF %s %s\r\n", transactionID, strings.Join(supportedAuthMethods, " "))))
}
