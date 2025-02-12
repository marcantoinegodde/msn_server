package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strings"

	"gorm.io/gorm"
)

func HandleBLP(db *gorm.DB, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	blp := database.Blp(args)
	switch blp {
	case database.AL, database.BL:
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

	if user.Blp == blp {
		SendError(c, tid, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Blp = blp
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	HandleSendBLP(c, tid, user.DataVersion, user.Blp)
	return nil
}

func HandleSendBLP(c *clients.Client, tid uint32, version uint32, blp database.Blp) {
	res := fmt.Sprintf("BLP %d %d %s\r\n", tid, version, blp)
	c.Send(res)
}
