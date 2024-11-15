package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func HandleSYN(conn net.Conn, db *gorm.DB, s *Session, arguments string) error {
	arguments, _, _ = strings.Cut(arguments, "\r\n")
	transactionID, arguments, err := parseTransactionID(arguments)
	if err != nil {
		return err
	}

	version, err := strconv.Atoi(arguments)
	if err != nil {
		return err
	}

	if !s.connected {
		SendError(conn, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", s.email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("SYN %s %d\r\n", transactionID, user.DataVersion)
	log.Println(">>>", res)
	conn.Write([]byte(res))

	if uint32(version) != user.DataVersion {
		// Start user's data synchronization

		// Send GTC
		HandleSendGTC(conn, transactionID, user.DataVersion, user.Gtc)

		// Send BLP
		HandleSendBLP(conn, transactionID, user.DataVersion, user.Blp)

		// Send LST FL
		HandleSendLST(conn, transactionID, "FL", &user)
		// Send LST AL
		HandleSendLST(conn, transactionID, "AL", &user)
		// Send LST BL
		HandleSendLST(conn, transactionID, "BL", &user)
		// Send LST RL
		HandleSendLST(conn, transactionID, "RL", &user)
	}

	return nil
}
