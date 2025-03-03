package commands

import (
	"errors"
	"msnserver/pkg/database"
	"strconv"
	"strings"
)

type cki struct {
	Cki       string `json:"cki"`
	SessionID uint32 `json:"session_id"`
}

func parseTransactionID(args string) (uint32, string, error) {
	rawTid, args, _ := strings.Cut(args, " ")

	parsedTid, err := strconv.ParseUint(rawTid, 10, 32)
	if err != nil {
		return 0, "", errors.New("invalid transaction ID")
	}
	tid := uint32(parsedTid)

	return tid, args, nil
}

func parseSessionID(sessionID string) (uint32, error) {
	parsedSessionID, err := strconv.ParseUint(sessionID, 10, 32)
	if err != nil {
		return 0, errors.New("invalid session ID")
	}

	return uint32(parsedSessionID), nil
}

func isMember(userList []*database.User, principal *database.User) bool {
	for _, u := range userList {
		if u.Email == principal.Email {
			return true
		}
	}

	return false
}
