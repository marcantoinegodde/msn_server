package commands

import "net"

func HandleOUT(conn net.Conn) {
	conn.Write([]byte("OUT\r\n"))
}
