package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

var supportedProtocols = []string{"MSNP2", "CVR0"}

func HandleVER(c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}
	clientProtocols := strings.Fields(args)

	serverProtocols := []string{}
	for _, protocol := range clientProtocols {
		for _, serverProtocol := range supportedProtocols {
			if protocol == serverProtocol {
				serverProtocols = append(serverProtocols, protocol)
			}
		}
	}

	if len(serverProtocols) < 2 {
		res := fmt.Sprintf("VER %d %s\r\n", tid, "0")
		c.Send(res)
		return errors.New("protocol mismatch")
	}

	res := fmt.Sprintf("VER %d %s\r\n", tid, strings.Join(serverProtocols, " "))
	c.Send(res)
	return nil
}
