package commands

import "msnserver/pkg/clients"

func HandleOUT(c *clients.Client) {
	res := "OUT\r\n"
	c.SendChan <- res
}
