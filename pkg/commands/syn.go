package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func HandleSYN(db *gorm.DB, c *clients.Client, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	tid, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	version, err := strconv.Atoi(arguments)
	if err != nil {
		return err
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("SYN %d %d\r\n", tid, user.DataVersion)
	c.Send(res)

	if uint32(version) != user.DataVersion {
		// Start user's data synchronization

		// Send GTC
		HandleSendGTC(c, tid, user.DataVersion, user.Gtc)

		// Send BLP
		HandleSendBLP(c, tid, user.DataVersion, user.Blp)

		// Send LST FL
		HandleSendLST(c, tid, "FL", &user)
		// Send LST AL
		HandleSendLST(c, tid, "AL", &user)
		// Send LST BL
		HandleSendLST(c, tid, "BL", &user)
		// Send LST RL
		HandleSendLST(c, tid, "RL", &user)
	}

	return nil
}
