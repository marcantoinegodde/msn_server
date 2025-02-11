package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

var supportedProtocols = []string{"MSNP2", "CVR0"}

func HandleVER(c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}
	clientProtocols := strings.Fields(arguments)

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
