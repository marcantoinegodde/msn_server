package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

const (
	recommendedVersion string = "1.0.0863"
	minimumVersion     string = "1.0.0863"
	downloadURL        string = "http://messenger.hotmail.com/mmsetup.exe"
	infoURL            string = "http://messenger.hotmail.com"
)

func HandleCVR(c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, _, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("CVR %s %s %s %s %s %s\r\n", transactionID, recommendedVersion, recommendedVersion, minimumVersion, downloadURL, infoURL)
	c.Send(res)
	return nil
}
