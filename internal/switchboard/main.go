package switchboard

import (
	"fmt"
	"log"
	"msnserver/config"
	"net"
)

type SwitchboardServer struct {
	config *config.MSNServerConfiguration
}

func NewSwitchboardServer(c *config.MSNServerConfiguration) *SwitchboardServer {
	return &SwitchboardServer{
		config: c,
	}
}

func (ss *SwitchboardServer) Start() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ss.config.SwitchboardServer.ServerAddr, ss.config.SwitchboardServer.ServerPort))
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

	defer ln.Close()

	log.Println("Listening on:", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go ss.handleConnection(conn)
	}
}

func (ss *SwitchboardServer) handleConnection(conn net.Conn) {
	defer conn.Close()
}
