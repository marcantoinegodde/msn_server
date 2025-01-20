package commands

import (
	"fmt"
	"strings"
)

func HandleSND(c chan string, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, _, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("SND %s OK\r\n", tid)
	c <- res
	return nil
}
