package notification

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/commands"
	"net"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type NotificationServer struct {
	db      *gorm.DB
	config  *config.MSNServerConfiguration
	m       sync.Mutex
	clients map[string]*Client
}

func NewNotificationServer(db *gorm.DB, c *config.MSNServerConfiguration) *NotificationServer {
	return &NotificationServer{
		db:      db,
		config:  c,
		m:       sync.Mutex{},
		clients: map[string]*Client{},
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
			continue
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go ns.handleConnection(conn)
	}
}

func (ns *NotificationServer) handleConnection(conn net.Conn) {
	c := &Client{
		id:       conn.RemoteAddr().String(),
		conn:     conn,
		sendChan: make(chan string),
	}

	defer func() {
		ns.m.Lock()
		delete(ns.clients, c.id)
		ns.m.Unlock()

		conn.Close()
		log.Println("Client disconnected:", conn.RemoteAddr())
	}()

	ns.m.Lock()
	ns.clients[c.id] = c
	ns.m.Unlock()

	go c.sendHandler()

	s := &commands.Session{}

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}

		go func() {
			data := string(buffer)
			log.Println("<<<", data)

			command, arguments, found := strings.Cut(data, " ")
			if !found {
				command, _, _ = strings.Cut(data, "\r\n")
			}

			switch command {
			case "VER":
				if err := commands.HandleVER(c.sendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "INF":
				if err := commands.HandleINF(c.sendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "USR":
				tid, err := commands.HandleReceiveUSR(s, arguments)
				if err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

				if err := commands.HandleSendUSR(c.sendChan, ns.db, s, tid); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "SYN":
				if err := commands.HandleSYN(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "CHG":
				if err := commands.HandleCHG(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "CVR":
				if err := commands.HandleCVR(c.sendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "GTC":
				if err := commands.HandleGTC(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "BLP":
				if err := commands.HandleBLP(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "ADD":
				if err := commands.HandleADD(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "REA":
				if err := commands.HandleREA(c.sendChan, ns.db, s, arguments); err != nil {
					log.Println("Error:", err)
					close(c.sendChan)
				}

			case "OUT":
				commands.HandleOUT(c.sendChan)
				close(c.sendChan)

			default:
				log.Println("Unknown command:", command)
				close(c.sendChan)
			}
		}()
	}
}
