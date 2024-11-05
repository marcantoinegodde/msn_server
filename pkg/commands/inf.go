package commands

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var supportedAuthMethods = []string{"MD5"}

func HandleINF(conn net.Conn, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, _, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("INF %s %s\r\n", transactionID, strings.Join(supportedAuthMethods, " "))
	log.Println(">>>", res)
	conn.Write([]byte(res))
	return nil
}
