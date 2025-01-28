package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

func HandleSND(c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, _, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c.SendChan, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	// We don't send any email, just ack the transaction
	res := fmt.Sprintf("SND %s OK\r\n", tid)
	c.SendChan <- res
	return nil
}
