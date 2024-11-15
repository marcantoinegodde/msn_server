package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net"
	"slices"
	"strings"

	"gorm.io/gorm"
)

func HandleReceiveUSR(conn net.Conn, db *gorm.DB, s *Session, arguments string) (string, error) {
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

	s.authMethod = splitArguments[0]
	s.authState = splitArguments[1]

	if splitArguments[1] == "I" {
		s.email = splitArguments[2]
	} else if splitArguments[1] == "S" {
		s.password = splitArguments[2]
	} else {
		err := errors.New("invalid auth state")
		return "", err
	}

	return transactionID, nil
}

func HandleSendUSR(conn net.Conn, db *gorm.DB, s *Session, transactionID string) error {
	switch s.authMethod {
	case "MD5":
		if s.authState == "I" {
			var user database.User
			query := db.First(&user, "email = ?", s.email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("user not found")
			} else if query.Error != nil {
				return query.Error
			}

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, s.authMethod, "S", user.Salt)
			log.Println(">>>", res)
			conn.Write([]byte(res))
			return nil
		} else if s.authState == "S" {
			var user database.User
			query := db.First(&user, "email = ?", s.email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("user not found")
			} else if query.Error != nil {
				return query.Error
			}

			if user.Password != s.password {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("invalid password")
			}

			s.connected = true

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, "OK", user.Email, user.Name)
			log.Println(">>>", res)
			conn.Write([]byte(res))
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
