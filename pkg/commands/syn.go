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

func HandleSYN(c chan string, db *gorm.DB, s *clients.Session, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	if !s.Authenticated {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	version, err := strconv.Atoi(arguments)
	if err != nil {
		return err
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("SYN %s %d\r\n", transactionID, user.DataVersion)
	c <- res

	if uint32(version) != user.DataVersion {
		// Start user's data synchronization

		// Send GTC
		HandleSendGTC(c, transactionID, user.DataVersion, user.Gtc)

		// Send BLP
		HandleSendBLP(c, transactionID, user.DataVersion, user.Blp)

		// Send LST FL
		HandleSendLST(c, transactionID, "FL", &user)
		// Send LST AL
		HandleSendLST(c, transactionID, "AL", &user)
		// Send LST BL
		HandleSendLST(c, transactionID, "BL", &user)
		// Send LST RL
		HandleSendLST(c, transactionID, "RL", &user)
	}

	return nil
}
