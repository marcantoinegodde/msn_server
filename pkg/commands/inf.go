package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

var supportedAuthMethods = []string{"MD5"}

func HandleINF(c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, _, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("INF %d %s\r\n", tid, strings.Join(supportedAuthMethods, " "))
	c.Send(res)
	return nil
}
