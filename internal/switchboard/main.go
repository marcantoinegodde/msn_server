package switchboard

import (
	"fmt"
	"log"
	"msnserver/config"
	"msnserver/pkg/clients"
	"msnserver/pkg/commands"
	"net"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SwitchboardServer struct {
	config  *config.MSNServerConfiguration
	db      *gorm.DB
	rdb     *redis.Client
	m       *sync.Mutex
	clients map[string]*clients.Client
}

func NewSwitchboardServer(c *config.MSNServerConfiguration, db *gorm.DB, rdb *redis.Client) *SwitchboardServer {
	return &SwitchboardServer{
		config:  c,
		db:      db,
		rdb:     rdb,
		m:       &sync.Mutex{},
		clients: map[string]*clients.Client{},
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
	c := clients.NewClient(conn)

	defer func() {
		if c.Session.Email != "" {
			ss.m.Lock()
			delete(ss.clients, c.Session.Email)
			ss.m.Unlock()
		}

		c.Disconnect()
	}()

	for {
		select {
		case msg := <-c.RecvChan:
			command, arguments, found := strings.Cut(msg, " ")
			if !found {
				command, _, _ = strings.Cut(msg, "\r\n")
			}

			switch command {
			case "USR":
				if err := commands.HandleUSRSwitchboard(ss.db, ss.rdb, ss.m, ss.clients, c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			default:
				log.Println("Unknown command:", command)
				return
			}

		case <-c.DoneChan:
			return
		}
	}
}
