package commands

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func HandleSYN(conn net.Conn, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	version, err := strconv.Atoi(arguments)
	if err != nil {
		return err
	}

	// TODO: Implement proper settings synchronization
	res := fmt.Sprintf("SYN %s %d\r\n", transactionID, version)
	log.Println(">>>", res)
	conn.Write([]byte(res))
	return nil
}
