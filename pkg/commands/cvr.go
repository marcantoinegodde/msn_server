package commands

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	RecommendedVersion = "1.0.0863"
	MinimumVersion     = "1.0.0863"
	DownloadURL        = "http://messenger.hotmail.com/mmsetup.exe"
	InfoURL            = "http://messenger.hotmail.com"
)

func HandleCVR(conn net.Conn, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	res := fmt.Sprintf("CVR %s %s %s %s %s %s\r\n", transactionID, RecommendedVersion, RecommendedVersion, MinimumVersion, DownloadURL, InfoURL)
	log.Println(">>>", res)
	conn.Write([]byte(res))
	return nil
}
