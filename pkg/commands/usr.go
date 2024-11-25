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

func HandleReceiveUSR(s *clients.Session, arguments string) (string, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return "", err
	}

	splitArguments := strings.Split(arguments, " ")
	if len(splitArguments) != 3 {
		err := errors.New("invalid transaction")
		return "", err
	}

	if !slices.Contains(supportedAuthMethods, splitArguments[0]) {
		err := errors.New("unsupported authentication method")
		return "", err
	}

	s.AuthMethod = splitArguments[0]
	s.AuthState = splitArguments[1]

	if splitArguments[1] == "I" {
		s.Email = splitArguments[2]
	} else if splitArguments[1] == "S" {
		s.Password = splitArguments[2]
	} else {
		err := errors.New("invalid auth state")
		return "", err
	}

	return transactionID, nil
}

func HandleSendUSR(c chan string, db *gorm.DB, s *clients.Session, transactionID string) error {
	switch s.AuthMethod {
	case "MD5":
		if s.AuthState == "I" {
			var user database.User
			query := db.First(&user, "email = ?", s.Email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(c, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("user not found")
			} else if query.Error != nil {
				return query.Error
			}

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, s.AuthMethod, "S", user.Salt)
			c <- res
			return nil

		} else if s.AuthState == "S" {
			var user database.User
			query := db.First(&user, "email = ?", s.Email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(c, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("user not found")
			} else if query.Error != nil {
				return query.Error
			}

			if user.Password != s.Password {
				SendError(c, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("invalid password")
			}

			s.Connected = true

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, "OK", user.Email, user.Name)
			c <- res
			return nil

		} else {
			err := errors.New("invalid auth state")
			return err
		}

	default:
		err := errors.New("unsupported authentication method")
		return err
	}
}
