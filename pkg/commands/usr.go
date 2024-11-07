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

type authParams struct {
	authMethod string
	authState  string
	email      string
	password   string
}

func HandleReceiveUSR(conn net.Conn, db *gorm.DB, arguments string) (string, authParams, error) {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return "", authParams{}, err
	}

	splitArguments := strings.Split(arguments, " ")
	if len(splitArguments) != 3 {
		err := errors.New("Invalid transaction")
		return "", authParams{}, err
	}

	authParams := authParams{authMethod: splitArguments[0], authState: splitArguments[1], email: splitArguments[2]}

	if !slices.Contains(supportedAuthMethods, authParams.authMethod) {
		err := errors.New("Unsupported authentication method")
		return "", authParams, err
	}

	return transactionID, authParams, nil
}

func HandleSendUSR(conn net.Conn, db *gorm.DB, transactionID string, authParams authParams) error {
	switch authParams.authMethod {
	case "MD5":
		if authParams.authState == "I" {
			var user database.User
			query := db.First(&user, "email = ?", authParams.email)
			if errors.Is(query.Error, gorm.ErrRecordNotFound) {
				SendError(conn, transactionID, ERR_AUTH_FAILED)
				return errors.New("User not found")
			} else if query.Error != nil {
				return query.Error
			}

			res := fmt.Sprintf("USR %s %s %s %s\r\n", transactionID, authParams.authMethod, "S", user.Salt)
			log.Println(">>>", res)
			conn.Write([]byte(res))
			return nil
		} else if authParams.authState == "S" {
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
