package commands

import (
	"errors"
	"msnserver/pkg/database"
	"strconv"
	"strings"
)

func parseTransactionID(arguments string) (uint32, string, error) {
	rawTid, arguments, _ := strings.Cut(arguments, " ")

	parsedTid, err := strconv.ParseUint(rawTid, 10, 32)
	if err != nil {
		return 0, "", errors.New("invalid transaction ID")
	}
	tid := uint32(parsedTid)

	return tid, arguments, nil
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
