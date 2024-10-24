package dispatch

import (
	"fmt"
	"msnserver/config"
	"msnserver/internal/commands"
	"net"
	"strings"
)

func StartDispatchServer() {
	fmt.Println("Starting MSN dispatch server...")

	config.LoadConfig()

	ln, err := net.Listen("tcp", config.Config.DispatchServer.ServerAddr+":"+config.Config.DispatchServer.ServerPort)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Listening on:", config.Config.DispatchServer.ServerAddr+":"+config.Config.DispatchServer.ServerPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for conn != nil {

		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		data := string(buffer)

		command, arguments, _ := strings.Cut(data, " ")

		switch command {
		case "VER":
			commands.HandleVER(conn, arguments)
		case "INF":
			commands.HandleINF(conn, transactionID)
		case "USR":
			commands.HandleUSRDispatch(conn, transactionID)
		case "OUT":
			commands.HandleOUT(conn)
		default:
			fmt.Println("Unknown command: ", command)
		}
	}

	fmt.Println("Client disconnected")
}
