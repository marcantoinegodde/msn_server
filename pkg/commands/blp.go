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

var blpMode = []string{"AL", "BL"}

func HandleBLP(c chan string, db *gorm.DB, s *clients.Session, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !s.Authenticated {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	if !slices.Contains(blpMode, args) {
		return errors.New("invalid mode")
	}

	var user database.User
	query := db.First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if user.Blp == args {
		SendError(c, transactionID, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Blp = args
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	HandleSendBLP(c, transactionID, user.DataVersion, user.Blp)
	return nil
}

func HandleSendBLP(c chan string, tid string, version uint32, blp string) {
	res := fmt.Sprintf("BLP %s %d %s\r\n", tid, version, blp)
	c <- res
}
