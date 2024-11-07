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

type AuthParams struct {
	authMethod string
	authState  string
	email      string
	password   string
	connected  bool
}

func HandleReceiveUSR(conn net.Conn, db *gorm.DB, ap *AuthParams, arguments string) (string, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return "", err
	}

	splitArguments := strings.Split(arguments, " ")
	if len(splitArguments) != 3 {
		err := errors.New("Invalid transaction")
		return "", err
	}

	if !slices.Contains(supportedAuthMethods, splitArguments[0]) {
		err := errors.New("Unsupported authentication method")
		return "", err
	}

	ap.authMethod = splitArguments[0]
	ap.authState = splitArguments[1]

	if splitArguments[1] == "I" {
		ap.email = splitArguments[2]
	} else if splitArguments[1] == "S" {
		ap.password = splitArguments[2]
	} else {
		err := errors.New("Invalid auth state")
		return "", err
	}

	return transactionID, nil
}

func HandleSendUSR(conn net.Conn, db *gorm.DB, ap *AuthParams, transactionID string) error {
	switch ap.authMethod {
	case "MD5":
		if ap.authState == "I" {
			var user database.User
			query := db.First(&user, "email = ?", ap.email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("User not found")
			} else if query.Error != nil {
				return query.Error
			}

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, ap.authMethod, "S", user.Salt)
			log.Println(">>>", res)
			conn.Write([]byte(res))
			return nil
		} else if ap.authState == "S" {
			var user database.User
			query := db.First(&user, "email = ?", ap.email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("User not found")
			} else if query.Error != nil {
				return query.Error
			}

			if user.Password != ap.password {
				SendError(conn, transactionID, ERR_AUTHENTICATION_FAILED)
				return errors.New("Invalid password")
			}

			ap.connected = true

			res := fmt.Sprintf("USR %s %s %s %s %d\r\n", transactionID, "OK", user.Email, user.Name, boolToInt(user.Verified))
			log.Println(">>>", res)
			conn.Write([]byte(res))
			return nil
		} else {
			err := errors.New("Invalid auth state")
			return err
		}

	default:
		err := errors.New("Unsupported authentication method")
		return err
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
