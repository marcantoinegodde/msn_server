package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"slices"
	"strings"
	"sync"

	"gorm.io/gorm"
)

func HandleUSR(db *gorm.DB, m *sync.Mutex, clients map[string]*clients.Client, c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	// Reject already authenticated clients
	if c.Session.Authenticated {
		SendError(c.SendChan, transactionID, ERR_ALREADY_LOGIN)
		return nil
	}

	// Parse arguments
	splitArguments := strings.Fields(arguments)
	if len(splitArguments) != 3 {
		err := errors.New("invalid transaction")
		return err
	}

	var authMethod = splitArguments[0]
	var authState = splitArguments[1]
	var password string

	// Validate authentication method
	if !slices.Contains(supportedAuthMethods, authMethod) {
		err := errors.New("unsupported authentication method")
		return err
	}

	switch authState {
	case "I":
		c.Session.Email = splitArguments[2]

	case "S":
		password = splitArguments[2]

	default:
		err := errors.New("invalid auth state")
		return err
	}

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c.SendChan, transactionID, ERR_AUTHENTICATION_FAILED)
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	switch authState {
	case "I":
		res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, authMethod, "S", user.Salt)
		c.SendChan <- res
		return nil

	case "S":
		if user.Password != password {
			SendError(c.SendChan, transactionID, ERR_AUTHENTICATION_FAILED)
			return errors.New("invalid password")
		}

		// Update user status
		c.Session.Authenticated = true

		// Update client map, handle logout if user is already logged in
		m.Lock()
		if oldClient, ok := clients[c.Session.Email]; ok {
			HandleOUT(oldClient, "OTH")
			oldClient.DoneChan <- true
			m.Unlock()
			oldClient.Wg.Wait()
			m.Lock()
		}
		clients[c.Session.Email] = c
		m.Unlock()

		res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, "OK", user.Email, user.DisplayName)
		c.SendChan <- res
		return nil

	default:
		err := errors.New("invalid auth state")
		return err
	}
}

/*
For now, we just return the tid without actually parsing the USR command
sent by the user. This could be improved later if we need to associate
a user account to a specific NS.
*/

func HandleUSRDispatch(arguments string) (string, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, _, err := parseTransactionID(arguments)
	if err != nil {
		return "", err
	}

	return tid, nil
}
