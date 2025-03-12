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

func HandleSYN(db *gorm.DB, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 1 {
		return errors.New("invalid transaction")
	}

	parsedVersion, err := strconv.ParseUint(splitArguments[0], 10, 32)
	if err != nil {
		return err
	}
	version := uint32(parsedVersion)

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("SYN %d %d\r\n", tid, user.DataVersion)
	c.Send(res)

	if version != user.DataVersion {
		// Start user's data synchronization

		// Send GTC
		HandleSendGTC(c, tid, user.DataVersion, user.Gtc)

		// Send BLP
		HandleSendBLP(c, tid, user.DataVersion, user.Blp)

		// Send LST FL
		HandleSendLST(c, tid, ForwardList, &user)
		// Send LST AL
		HandleSendLST(c, tid, AllowList, &user)
		// Send LST BL
		HandleSendLST(c, tid, BlockList, &user)
		// Send LST RL
		HandleSendLST(c, tid, ReverseList, &user)
	}

	return nil
}
