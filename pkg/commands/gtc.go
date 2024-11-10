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

var gtcMode = []string{"A", "N"}

func HandleGTC(conn net.Conn, db *gorm.DB, ap *AuthParams, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	var user database.User
	query := db.First(&user, "email = ?", ap.email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if !slices.Contains(gtcMode, args) {
		return errors.New("invalid mode")
	}

	if user.Gtc == args {
		SendError(conn, transactionID, ERR_ALREADY_IN_THE_MODE)
		return errors.New("user already in requested mode")
	}

	user.Gtc = args
	user.DataVersion++
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("GTC %s %d %s\r\n", transactionID, user.DataVersion, user.Gtc)
	log.Println(">>>", res)
	conn.Write([]byte(res))
	return nil
}
