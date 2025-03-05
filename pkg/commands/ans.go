package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"msnserver/pkg/sessions"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func HandleANS(db *gorm.DB, rdb *redis.Client, sbs *sessions.SwitchboardSessions, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	// Reject already authenticated clients
	if c.Session.Authenticated {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("already authenticated")
	}

	// Parse arguments
	splitArguments := strings.Fields(args)
	if len(splitArguments) != 3 {
		return errors.New("invalid transaction")
	}

	email := splitArguments[0]
	userCki := splitArguments[1]
	sessionID := splitArguments[2]

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

	// Parse session ID
	sid, err := parseSessionID(sessionID)
	if err != nil {
		return err
	}

	// Validate session ID
	if sid != cki.SessionID {
		SendError(c, tid, ERR_AUTHENTICATION_FAILED)
		return errors.New("invalid session ID")
	}

	// Fetch user from the database
	var user database.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		return err
	}

	// Join session
	s, err := sbs.JoinSession(c, sid)
	if err != nil {
		return err
	}

	// Update client session
	c.Session.Email = user.Email
	c.Session.DisplayName = user.DisplayName
	c.Session.Authenticated = true

	// Send initial roaster information to the client
	// Even if clients disconnect in the meantime, pointers to the clients are still valid
	HandleSendIRO(c, tid, s)

	// Send join notification to all clients in the session
	HandleSendJOI(c, s)

	res := fmt.Sprintf("ANS %d OK\r\n", tid)
	c.Send(res)

	return nil
}
