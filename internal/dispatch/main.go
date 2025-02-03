package dispatch

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/clients"
	"msnserver/pkg/commands"
	"net"
	"strings"
)

type DispatchServer struct {
	config *config.MSNServerConfiguration
}

func NewDispatchServer(c *config.MSNServerConfiguration) *DispatchServer {
	return &DispatchServer{
		config: c,
	}
}

func (ds *DispatchServer) Start() {
	ln, err := net.Listen("tcp", ds.config.DispatchServer.ServerAddr+":"+ds.config.DispatchServer.ServerPort)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

	defer ln.Close()

	log.Println("Listening on:", ds.config.DispatchServer.ServerAddr+":"+ds.config.DispatchServer.ServerPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go ds.handleConnection(conn)
	}
}

func (ds *DispatchServer) handleConnection(conn net.Conn) {
	c := clients.NewClient(conn)

	defer c.Disconnect()

	for {
		select {
		case msg := <-c.RecvChan:
			command, arguments, found := strings.Cut(msg, " ")
			if !found {
				command, _, _ = strings.Cut(msg, "\r\n")
			}

			switch command {
			case "VER":
				if err := commands.HandleVER(c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "INF":
				if err := commands.HandleINF(c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "USR":
				tid, err := commands.HandleUSRDispatch(arguments)
				if err != nil {
					log.Println("Error:", err)
					return
				}

				commands.HandleXFR(ds.config.DispatchServer, c, tid)
				return

			case "OUT":
				commands.HandleOUT(c, "")
				return

			default:
				log.Println("Unknown command:", command)
				return
			}

		case <-c.DoneChan:
			return
		}
	}
}
