package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

func HandleCAL(db *gorm.DB, c *clients.Client, args string) error {
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

	calleeEmail := splitArguments[0]

	// Fetch the user from the database
	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	// Validate the callee's email
	if !utils.IsValidEmail(calleeEmail) {
		SendError(c, tid, ERR_INVALID_USERNAME)
		return nil
	}

	// Fetch the callee from the database
	var callee database.User
	query = db.Preload("AllowList").Preload("BlockList").First(&callee, "email = ?", calleeEmail)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c, tid, ERR_NOT_ONLINE)
		return nil
	} else if query.Error != nil {
		return query.Error
	}

	// Check if the callee is online
	if callee.Status == database.FLN || callee.Status == database.HDN {
		SendError(c, tid, ERR_NOT_ONLINE)
		return nil
	}

	// Check if the caller isn't blocked by the callee
	if isMember(callee.BlockList, &user) {
		SendError(c, tid, ERR_NOT_ONLINE)
		return nil
	}

	// Check if the caller is allowed to call the callee if the callee is in BL mode
	if callee.Blp == database.BL && !isMember(callee.AllowList, &user) {
		SendError(c, tid, ERR_NOT_ONLINE)
		return nil
	}

	// TODO: Create a session
	// TODO: Generate a random session ID
	// TODO: Handle self adding
	// TODO: Max session
	// TODO: Communicate invitation to callee
	sessionId := 12345678
	res := fmt.Sprintf("CAL %d RINGING %d\r\n", tid, sessionId)
	c.Send(res)

	return nil
}
