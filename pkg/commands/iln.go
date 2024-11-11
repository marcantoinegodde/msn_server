package commands

import (
	"fmt"
	"log"
	"net"
)

func HandleSendILN(conn net.Conn, transactionID string, status string, email string, name string) {
	res := fmt.Sprintf("ILN %s %s %s %s\r\n", transactionID, status, email, name)
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
