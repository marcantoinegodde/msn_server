package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var blpMode = []string{"AL", "BL"}

func HandleBLP(conn net.Conn, db *gorm.DB, s *Session, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	var user database.User
	query := db.First(&user, "email = ?", s.email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if !slices.Contains(blpMode, args) {
		return errors.New("invalid mode")
	}

	if user.Blp == args {
		SendError(conn, transactionID, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Blp = args
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("BLP %s %d %s\r\n", transactionID, user.DataVersion, user.Blp)
	log.Println(">>>", res)
	conn.Write([]byte(res))
	return nil
}
