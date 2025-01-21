package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

func HandleSND(c chan string, s *clients.Session, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, _, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !s.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	// We don't send any email, just ack the transaction
	res := fmt.Sprintf("SND %s OK\r\n", tid)
	c <- res
	return nil
}
