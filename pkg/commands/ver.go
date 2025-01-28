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
	transactionID, arguments, err := parseTransactionID(arguments)
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
		res := fmt.Sprintf("VER %s %s\r\n", transactionID, "0")
		c.SendChan <- res
		return errors.New("protocol mismatch")
	}

	res := fmt.Sprintf("VER %s %s\r\n", transactionID, strings.Join(serverProtocols, " "))
	c.SendChan <- res
	return nil
}
