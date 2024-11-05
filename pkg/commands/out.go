package commands

import (
	"log"
	"net"
)

func HandleOUT(conn net.Conn) {
	res := "OUT\r\n"
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
