package dispatch

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/commands"
	"net"
	"strings"

	"gorm.io/gorm"
)

type DispatchServer struct {
	db     *gorm.DB
	config *config.MSNServerConfiguration
}

func NewDispatchServer(db *gorm.DB, c *config.MSNServerConfiguration) *DispatchServer {
	return &DispatchServer{
		db:     db,
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
			return
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go ds.handleConnection(conn)
	}
}

func (ds *DispatchServer) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		} else {
			log.Println("Client disconnected:", conn.RemoteAddr())
		}
	}()

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		s := &commands.Session{}

		data := string(buffer)
		log.Println("<<<", data)

		command, arguments, found := strings.Cut(data, " ")
		if !found {
			command, _, _ = strings.Cut(data, "\r\n")
		}

		switch command {
		case "VER":
			err := commands.HandleVER(conn, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "INF":
			err := commands.HandleINF(conn, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "USR":
			transactionID, err := commands.HandleReceiveUSR(conn, ds.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			commands.HandleXFR(conn, ds.config.DispatchServer, transactionID)
			return

		case "OUT":
			commands.HandleOUT(conn)
			return

		default:
			log.Println("Unknown command:", command)
			return
		}
	}
}
