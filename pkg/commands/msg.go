package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/sessions"
	"strconv"
	"strings"
)

func HandleMSG(sbs *sessions.SwitchboardSessions, c *clients.Client, args string) error {
	args, msg, _ := strings.Cut(args, "\r\n")
	tid, arguments, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	// Parse arguments
	splitArguments := strings.Fields(arguments)
	if len(splitArguments) != 2 {
		err := errors.New("invalid transaction")
		return err
	}

	ackMode := splitArguments[0]
	rawLength := splitArguments[1]

	length, err := strconv.Atoi(rawLength)
	if err != nil {
		return err
	}

	// Check message length
	if length > 462 {
		err := errors.New("message too long")
		return err
	}

	// Send the message to all clients in the session
	res := fmt.Sprintf("MSG %s %s %d\r\n%s", c.Session.Email, c.Session.DisplayName, length, msg[:length])
	if err := sbs.MessageSession(c, res); err != nil {
		if ackMode == "A" || ackMode == "N" {
			HandleSendNAK(c, tid)
		}
		log.Println("Error:", err)
		return nil
	}

	// Send ACK if ack mode is A
	if ackMode == "A" {
		HandleSendACK(c, tid)
	}

	return nil
}
