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

var gtcMode = []string{"A", "N"}

func HandleGTC(db *gorm.DB, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	if !slices.Contains(gtcMode, args) {
		return errors.New("invalid mode")
	}

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if user.Gtc == args {
		SendError(c, tid, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Gtc = args
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	HandleSendGTC(c, tid, user.DataVersion, user.Gtc)
	return nil
}

func HandleSendGTC(c *clients.Client, tid uint32, version uint32, gtc string) {
	res := fmt.Sprintf("GTC %d %d %s\r\n", tid, version, gtc)
	c.Send(res)
}
