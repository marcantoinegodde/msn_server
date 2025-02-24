package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/sessions"
	"strings"

	"github.com/redis/go-redis/v9"
)

func HandleANS(rdb *redis.Client, sbs *sessions.SwitchboardSessions, c *clients.Client, args string) error {
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

	c.Session.Email = splitArguments[0]
	userCki := splitArguments[1]
	sessionID := splitArguments[2]

	// Fetch CKI from Redis
	rawCki, err := rdb.GetDel(context.TODO(), c.Session.Email).Result()
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

	// Join session
	if err := sbs.JoinSession(c, sid); err != nil {
		return err
	}

	c.Session.Authenticated = true

	res := fmt.Sprintf("ANS %d OK\r\n", tid)
	c.Send(res)

	return nil
}
