package notification

import (
	"log"
	"msnserver/config"
	"msnserver/pkg/clients"
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
	clients map[string]*clients.Client
}

func NewNotificationServer(db *gorm.DB, c *config.MSNServerConfiguration) *NotificationServer {
	return &NotificationServer{
		db:      db,
		config:  c,
		m:       sync.Mutex{},
		clients: map[string]*clients.Client{},
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
	c := &clients.Client{
		Id:       conn.RemoteAddr().String(),
		Conn:     conn,
		SendChan: make(chan string),
		Session:  &clients.Session{},
	}

	defer func() {
		ns.m.Lock()
		delete(ns.clients, c.Session.Email)
		ns.m.Unlock()

		conn.Close()
		log.Println("Client disconnected:", conn.RemoteAddr())
	}()

	go c.SendHandler()

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from connection:", err)
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
				if err := commands.HandleVER(c.SendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "INF":
				if err := commands.HandleINF(c.SendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "USR":
				tid, err := commands.HandleReceiveUSR(c.Session, arguments)
				if err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

				if err := commands.HandleSendUSR(c.SendChan, ns.db, c.Session, tid); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

				ns.m.Lock()
				ns.clients[c.Session.Email] = c
				ns.m.Unlock()

			case "SYN":
				if err := commands.HandleSYN(c.SendChan, ns.db, c.Session, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "CHG":
				if err := commands.HandleCHG(c.SendChan, ns.db, c.Session, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "CVR":
				if err := commands.HandleCVR(c.SendChan, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "GTC":
				if err := commands.HandleGTC(c.SendChan, ns.db, c.Session, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "BLP":
				if err := commands.HandleBLP(c.SendChan, ns.db, c.Session, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "ADD":
				if err := commands.HandleADD(c.SendChan, ns.db, c.Session, ns.clients, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "REA":
				if err := commands.HandleREA(c.SendChan, ns.db, c.Session, arguments); err != nil {
					log.Println("Error:", err)
					close(c.SendChan)
				}

			case "OUT":
				commands.HandleOUT(c.SendChan)
				close(c.SendChan)

			default:
				log.Println("Unknown command:", command)
				close(c.SendChan)
			}
		}()
	}
}
