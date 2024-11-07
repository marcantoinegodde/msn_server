package dispatch

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/commands"
	"net"
	"strings"

	"gorm.io/gorm"
)

func StartDispatchServer(db *gorm.DB) {
	ln, err := net.Listen("tcp", config.Config.DispatchServer.ServerAddr+":"+config.Config.DispatchServer.ServerPort)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

	defer ln.Close()

	log.Println("Listening on:", config.Config.DispatchServer.ServerAddr+":"+config.Config.DispatchServer.ServerPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			return
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *gorm.DB) {
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

		ap := &commands.AuthParams{}

		data := string(buffer)
		log.Println("<<<", data)

		command, arguments, _ := strings.Cut(data, " ")

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
			transactionID, err := commands.HandleReceiveUSR(conn, db, arguments, ap)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			commands.HandleXFR(conn, transactionID)
			return

		case "OUT":
			commands.HandleOUT(conn)

		default:
			log.Println("Unknown command:", command)
			return
		}
	}
}
