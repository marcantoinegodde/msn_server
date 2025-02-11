package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"strings"
)

/*
The URL command provides the client with a URL to connect to a specific service.

The received syntax is: URL <transaction-id> <url-type> [<parameter>]
The server will respond with: URL <transaction-id> <redirect-url> <authentication-url>

When receiving an anwser from the server, the client will create a temporary HTML file
with the necessary content to redirect the user onto the authentication page.

It then triggers a POST request to the authentication URL:
- Content-Type: application/x-www-form-urlencoded
- Parameters: login=<email>&k1=<k1>&k2=<k2>&rru=<redirect-url>&js=yes
*/

var urlArgs = map[string]string{
	"PASSWORD": "/password http://127.0.0.1",
	"PERSON":   "/person http://127.0.0.1",
	"COMPOSE":  "/compose http://127.0.0.1",
	"INBOX":    "/inbox http://127.0.1",
	"FOLDERS":  "/inbox http://127.0.0.1",
	"MESSAGE":  "/message http://127.0.0.1"}

func HandleURL(c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	var urlType string
	// var parameter string

	splitArguments := strings.Fields(args)
	if len(splitArguments) == 1 {
		urlType = splitArguments[0]
	} else if len(splitArguments) == 2 {
		urlType = splitArguments[0]
		// parameter = splitArguments[1]
	} else {
		return errors.New("invalid transaction")
	}

	url, ok := urlArgs[urlType]
	if !ok {
		SendError(c, tid, ERR_INVALID_PARAMETER)
		return fmt.Errorf("invalid URL type: %s", urlType)
	}

	res := fmt.Sprintf("URL %d %s\r\n", tid, url)
	c.Send(res)
	return nil
}
