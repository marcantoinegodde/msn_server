package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"msnserver/pkg/sessions"
	"slices"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SequenceValue string

const (
	Initiate   SequenceValue = "I"
	Subsequent SequenceValue = "S"
)

/*
For now, we just return the tid without actually parsing the USR command
sent by the user. This could be improved later if we need to associate
a user account to a specific NS.
*/

func HandleUSRDispatch(arguments string) (uint32, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, _, err := parseTransactionID(arguments)
	if err != nil {
		return 0, err
	}

	return tid, nil
}

func HandleUSR(db *gorm.DB, m *sync.Mutex, clients map[string]*clients.Client, c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	// Reject already authenticated clients
	if c.Session.Authenticated {
		SendError(c, tid, ERR_ALREADY_LOGIN)
		return nil
	}

	// Parse arguments
	splitArguments := strings.Fields(arguments)
	if len(splitArguments) != 3 {
		err := errors.New("invalid transaction")
		return err
	}

	var sp = splitArguments[0]
	var sv = SequenceValue(splitArguments[1])
	var password string

	// Validate security package
	if !slices.Contains(supportedSecurityPackages, sp) {
		err := errors.New("unsupported security package")
		return err
	}

	switch sv {
	case Initiate:
		c.Session.Email = splitArguments[2]

	case Subsequent:
		password = splitArguments[2]

	default:
		err := errors.New("invalid sequence value")
		return err
	}

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	switch sv {
	case Initiate:
		res := fmt.Sprintf("USR %d %s %s %s\r\n", tid, sp, Subsequent, user.Salt)
		c.Send(res)
		return nil

	case Subsequent:
		if user.Password != password {
			SendError(c, tid, ERR_AUTHENTICATION_FAILED)
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

		res := fmt.Sprintf("USR %d %s %s %s\r\n", tid, "OK", user.Email, user.DisplayName)
		c.Send(res)
		return nil

	default:
		err := errors.New("invalid sequence value")
		return err
	}
}

func HandleUSRSwitchboard(db *gorm.DB, rdb *redis.Client, sbs *sessions.SwitchboardSessions, c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	// Reject already authenticated clients
	if c.Session.Authenticated {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("already authenticated")
	}

	// Parse arguments
	splitArguments := strings.Fields(arguments)
	if len(splitArguments) != 2 {
		err := errors.New("invalid transaction")
		return err
	}

	email := splitArguments[0]
	userCki := splitArguments[1]

	// Fetch CKI from Redis
	rawCki, err := rdb.GetDel(context.TODO(), email).Result()
	if err == redis.Nil {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("cki not found")
	} else if err != nil {
		return errors.New("error getting cki")
	}

	var cki cki
	if err := json.Unmarshal([]byte(rawCki), &cki); err != nil {
		return err
	}

	// Validate CKI
	if userCki != cki.Cki {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("invalid cki")
	}

	var user database.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return err
	}

	// Create Switchboard session
	sbs.CreateSession(c)

	// Update client session
	c.Session.Email = user.Email
	c.Session.DisplayName = user.DisplayName
	c.Session.Authenticated = true

	res := fmt.Sprintf("USR %d OK %s %s\r\n", tid, user.Email, user.DisplayName)
	c.Send(res)

	return nil
}
