package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

var supportedSecurityPackages = []string{"MD5"}

func HandleINF(c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, _, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("INF %d %s\r\n", tid, strings.Join(supportedSecurityPackages, " "))
	c.Send(res)
	return nil
}
