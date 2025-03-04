package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

var supportedSecurityPackages = []string{"MD5"}

func HandleINF(c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, _, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("INF %d %s\r\n", tid, strings.Join(supportedSecurityPackages, " "))
	c.Send(res)
	return nil
}
