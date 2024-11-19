package notification

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/commands"
	"net"
	"strings"

	"gorm.io/gorm"
)

type NotificationServer struct {
	db     *gorm.DB
	config *config.MSNServerConfiguration
}

func NewNotificationServer(db *gorm.DB, c *config.MSNServerConfiguration) *NotificationServer {
	return &NotificationServer{
		db:     db,
		config: c,
	}
}

func (ns *NotificationServer) Start() {
	ln, err := net.Listen("tcp", ns.config.NotificationServer.ServerAddr+":"+ns.config.NotificationServer.ServerPort)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

	defer ln.Close()

	log.Println("Listening on:", ns.config.NotificationServer.ServerAddr+":"+ns.config.NotificationServer.ServerPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			return
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go ns.handleConnection(conn)
	}
}

func (ns *NotificationServer) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		} else {
			log.Println("Client disconnected:", conn.RemoteAddr())
		}
	}()

	s := &commands.Session{}

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error:", err)
			return
		}

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
			tid, err := commands.HandleReceiveUSR(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			err = commands.HandleSendUSR(conn, ns.db, s, tid)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "SYN":
			err := commands.HandleSYN(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "CHG":
			err := commands.HandleCHG(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "CVR":
			err := commands.HandleCVR(conn, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "GTC":
			err := commands.HandleGTC(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "BLP":
			err := commands.HandleBLP(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "ADD":
			err := commands.HandleADD(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "REA":
			err := commands.HandleREA(conn, ns.db, s, arguments)
			if err != nil {
				log.Println("Error:", err)
				return
			}

		case "OUT":
			commands.HandleOUT(conn)
			return

		default:
			log.Println("Unknown command:", command)
			return
		}
	}
}
