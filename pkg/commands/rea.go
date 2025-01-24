package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

var blockedWords = []string{"microsoft", "msn", "fuck"}

func HandleREA(c chan string, db *gorm.DB, s *clients.Session, clients map[string]*clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !s.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 2 {
		return errors.New("invalid transaction")
	}

	email, newDisplayName := splitArguments[0], splitArguments[1]

	if _, err := url.PathUnescape(newDisplayName); err != nil {
		return errors.New("invalid new name")
	}

	if strings.ContainsAny(newDisplayName, " \t\n\r") {
		return errors.New("invalid new name")
	}

	if len(newDisplayName) > 129 {
		return errors.New("new name too long")
	}

	for _, word := range blockedWords {
		if strings.Contains(strings.ToLower(newDisplayName), word) {
			SendError(c, tid, ERR_INVALID_FRIENDLY_NAME)
			return nil
		}
	}

	var user database.User
	query := db.First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if s.Email == email {
		user.DisplayName = newDisplayName
		user.DataVersion++
		query = db.Save(&user)
		if query.Error != nil {
			return query.Error
		}

		res := fmt.Sprintf("REA %s %d %s %s\r\n", tid, user.DataVersion, user.Email, user.DisplayName)
		c <- res

		if err := HandleBatchNLN(db, clients, s); err != nil {
			log.Println("Error:", err)
		}

	} else {
		// TODO: Improve this, store the nicknames

		user.DataVersion++
		query = db.Save(&user)
		if query.Error != nil {
			return query.Error
		}

		res := fmt.Sprintf("REA %s %d %s %s\r\n", tid, user.DataVersion, email, newDisplayName)
		c <- res
	}

	return nil
}
