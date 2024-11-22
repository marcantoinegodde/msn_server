package commands

import (
	"fmt"
	"strings"
)

const (
	recommendedVersion = "1.0.0863"
	minimumVersion     = "1.0.0863"
	downloadURL        = "http://messenger.hotmail.com/mmsetup.exe"
	infoURL            = "http://messenger.hotmail.com"
)

func HandleCVR(c chan string, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("CVR %s %s %s %s %s %s\r\n", transactionID, recommendedVersion, recommendedVersion, minimumVersion, downloadURL, infoURL)
	c <- res
	return nil
}
