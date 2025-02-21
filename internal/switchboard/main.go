package switchboard

import (
	"fmt"
	"log"
	"msnserver/config"
	"msnserver/pkg/clients"
	"msnserver/pkg/commands"
	"msnserver/pkg/sessions"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SwitchboardServer struct {
	config *config.MSNServerConfiguration
	db     *gorm.DB
	rdb    *redis.Client
	sbs    *sessions.SwitchboardSessions
}

func NewSwitchboardServer(c *config.MSNServerConfiguration, db *gorm.DB, rdb *redis.Client) *SwitchboardServer {
	return &SwitchboardServer{
		config: c,
		db:     db,
		rdb:    rdb,
		sbs:    sessions.NewSwitchboardSessions(),
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

	defer c.Disconnect()

	for {
		select {
		case msg := <-c.RecvChan:
			command, arguments, found := strings.Cut(msg, " ")
			if !found {
				command, _, _ = strings.Cut(msg, "\r\n")
			}

			switch command {
			case "USR":
				if err := commands.HandleUSRSwitchboard(ss.db, ss.rdb, ss.sbs, c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "CAL":
				if err := commands.HandleCAL(ss.config, ss.db, ss.rdb, ss.sbs, c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "ANS":
				if err := commands.HandleANS(ss.rdb, c, arguments); err != nil {
					log.Println("Error:", err)
					return
				}

			case "OUT":
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
