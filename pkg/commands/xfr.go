package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"msnserver/config"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	SB_SECURITY_PACKAGE string        = "CKI"
	CKI_TIMEOUT         time.Duration = 2 * time.Minute
)

func HandleXFRDispatch(cf config.DispatchServer, c *clients.Client, tid uint32) {
	res := fmt.Sprintf("XFR %d NS %s:%d\r\n", tid, cf.NotificationServerAddr, cf.NotificationServerPort)
	c.Send(res)
}

func HandleXFR(cf config.NotificationServer, db *gorm.DB, rdb *redis.Client, c *clients.Client, args string) error {
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

	referralType := splitArguments[0]

	if referralType != "SB" {
		SendError(c, tid, ERR_INVALID_PARAMETER)
		return errors.New("invalid parameter")
	}

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if query.Error != nil {
		return query.Error
	}

	if user.Status == database.HDN {
		SendError(c, tid, ERR_NOT_ALLOWED_WHEN_OFFLINE)
		return nil
	}

	cki := cki{
		Cki:       utils.GenerateRandomString(25),
		SessionID: 0,
	}

	jsonCki, err := json.Marshal(cki)
	if err != nil {
		return err
	}

	if err := rdb.Set(context.TODO(), c.Session.Email, jsonCki, CKI_TIMEOUT).Err(); err != nil {
		return err
	}

	res := fmt.Sprintf("XFR %d SB %s:%d %s %s\r\n", tid, cf.SwitchboardServerAddr, cf.SwitchboardServerPort, SB_SECURITY_PACKAGE, cki.Cki)
	c.Send(res)

	return nil
}
