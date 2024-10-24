package commands

import (
	"fmt"
	"net"
	"strings"
)

var supportedProtocols = []string{"MSNP2", "CVR0"}

func HandleVER(conn net.Conn, arguments string) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		conn.Close()
		return
	}
	clientProtocols := strings.Split(arguments, " ")

	serverProtocols := []string{}
	for _, protocol := range clientProtocols {
		for _, serverProtocol := range supportedProtocols {
			if protocol == serverProtocol {
				serverProtocols = append(serverProtocols, protocol)
			}
		}
	}

	if len(serverProtocols) < 2 {
		conn.Write([]byte(fmt.Sprintf("VER %s %s\r\n", transactionID, "0")))
		conn.Close()
		return
	}

	conn.Write([]byte(fmt.Sprintf("VER %s %s\r\n", transactionID, strings.Join(serverProtocols, " "))))
}
