package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

var blockedWords = []string{"microsoft", "msn", "fuck"}

func HandleREA(conn net.Conn, db *gorm.DB, s *Session, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 2 {
		return errors.New("invalid transaction")
	}

	email, newName := splitArguments[0], splitArguments[1]

	if _, err := url.PathUnescape(newName); err != nil {
		return errors.New("invalid new name")
	}

	if strings.ContainsAny(newName, " \t\n\r") {
		return errors.New("invalid new name")
	}

	if len(newName) > 129 {
		return errors.New("new name too long")
	}

	for _, word := range blockedWords {
		if strings.Contains(strings.ToLower(newName), word) {
			SendError(conn, tid, ERR_INVALID_FRIENDLY_NAME)
			return nil
		}
	}

	if !s.connected {
		SendError(conn, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	if s.email == email {
		// TODO: Add asynchronous communication to other users

		var user database.User
		query := db.First(&user, "email = ?", s.email)
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		} else if query.Error != nil {
			return query.Error
		}

		user.Name = newName
		user.DataVersion++
		query = db.Save(&user)
		if query.Error != nil {
			return query.Error
		}

		res := fmt.Sprintf("REA %s %d %s %s\r\n", tid, user.DataVersion, user.Email, user.Name)
		log.Println(">>>", res)
		conn.Write([]byte(res))

	} else {
		// TODO: Add principal's name modification
		log.Println("Not implemented")
	}

	return nil
}
