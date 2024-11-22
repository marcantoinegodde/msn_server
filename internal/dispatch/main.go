package dispatch

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/commands"
	"net"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type DispatchServer struct {
	db      *gorm.DB
	config  *config.MSNServerConfiguration
	m       sync.Mutex
	clients map[string]*Client
}

func NewDispatchServer(db *gorm.DB, c *config.MSNServerConfiguration) *DispatchServer {
	return &DispatchServer{
		db:      db,
		config:  c,
		m:       sync.Mutex{},
		clients: map[string]*Client{},
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
		ds.m.Lock()
		delete(ds.clients, conn.RemoteAddr().String())
		ds.m.Unlock()

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

		c := &Client{
			id:       conn.RemoteAddr().String(),
			conn:     conn,
			sendChan: make(chan string),
		}

		ds.m.Lock()
		ds.clients[c.id] = c
		ds.m.Unlock()

		go c.sendHandler()

		s := &commands.Session{}

		data := string(buffer)
		log.Println("<<<", data)

		command, arguments, found := strings.Cut(data, " ")
		if !found {
			command, _, _ = strings.Cut(data, "\r\n")
		}

		// TO FIX: Terrible code to be rewritten, async goroutine can't close the connection
		go func() {
			switch command {
			case "VER":
				if err := commands.HandleVER(c.sendChan, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "INF":
				if err := commands.HandleINF(c.sendChan, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "USR":
				transactionID, err := commands.HandleReceiveUSR(s, arguments)
				if err != nil {
					log.Println("Error:", err)
					return
				}

				commands.HandleXFR(c.sendChan, ds.config.DispatchServer, transactionID)
				return

			case "OUT":
				commands.HandleOUT(c.sendChan)
				return

			default:
				log.Println("Unknown command:", command)
				return
			}
		}()
	}
}
