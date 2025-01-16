package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func HandleUSR(c chan string, db *gorm.DB, s *clients.Session, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	splitArguments := strings.Split(arguments, " ")
	if len(splitArguments) != 3 {
		err := errors.New("invalid transaction")
		return err
	}

	s.Authenticated = false

	var authMethod = splitArguments[0]
	var authState = splitArguments[1]
	var password string

	if !slices.Contains(supportedAuthMethods, authMethod) {
		err := errors.New("unsupported authentication method")
		return err
	}

	switch authState {
	case "I":
		s.Email = splitArguments[2]

	case "S":
		password = splitArguments[2]

	default:
		err := errors.New("invalid auth state")
		return err
	}

	var user database.User
	query := db.First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c, transactionID, ERR_AUTHENTICATION_FAILED)
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	switch authState {
	case "I":
		res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, authMethod, "S", user.Salt)
		c <- res
		return nil

	case "S":
		if user.Password != password {
			SendError(c, transactionID, ERR_AUTHENTICATION_FAILED)
			return errors.New("invalid password")
		}

		s.Authenticated = true

		res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, "OK", user.Email, user.Name)
		c <- res
		return nil

	default:
		err := errors.New("invalid auth state")
		return err
	}
}

func HandleUSRDispatch(arguments string) (string, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, _, err := parseTransactionID(arguments)
	if err != nil {
		return "", err
	}

	return tid, nil
}
