package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strings"

	"gorm.io/gorm"
)

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

	gtc := database.Gtc(args)
	switch gtc {
	case database.A, database.N:
		break
	default:
		return errors.New("invalid mode")
	}

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if user.Gtc == gtc {
		SendError(c, tid, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Gtc = gtc
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	HandleSendGTC(c, tid, user.DataVersion, user.Gtc)
	return nil
}

func HandleSendGTC(c *clients.Client, tid uint32, version uint32, gtc database.Gtc) {
	res := fmt.Sprintf("GTC %d %d %s\r\n", tid, version, gtc)
	c.Send(res)
}
