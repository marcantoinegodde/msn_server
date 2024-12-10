package commands

import (
	"errors"
	"math"
	"msnserver/pkg/database"
	"strconv"
	"strings"
)

func parseTransactionID(arguments string) (string, string, error) {
	transactionID, arguments, _ := strings.Cut(arguments, " ")

	parsedTransactionID, err := strconv.Atoi(transactionID)
	if err != nil {
		return "", "", errors.New("invalid transaction ID")
	}

	if parsedTransactionID < 0 {
		return "", "", errors.New("invalid transaction ID")
	}
	if parsedTransactionID > (int(math.Pow(2, 32)) - 1) {
		return "", "", errors.New("invalid transaction ID")
	}

	return transactionID, arguments, nil
}

func isMember(userList []*database.User, principal *database.User) bool {
	for _, u := range userList {
		if u.Email == principal.Email {
			return true
		}
	}

	return false
}
