package commands

import (
	"fmt"
	"msnserver/pkg/clients"

	"golang.org/x/exp/slices"
)

var validReasons = []string{"OTH", "SSD"}

func HandleOUT(c *clients.Client, reason string) {
	var res string

	if slices.Contains(validReasons, reason) {
		res = fmt.Sprintf("OUT %s\r\n", reason)
	} else {
		res = "OUT\r\n"
	}

	c.Send(res)
}
