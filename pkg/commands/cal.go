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

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func HandleCAL(cf *config.MSNServerConfiguration, db *gorm.DB, rdb *redis.Client, c *clients.Client, args string) error {
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

	// Prevent the user from calling themselves
	if c.Session.Email == calleeEmail {
		SendError(c, tid, ERR_ALREADY_THERE)
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
	// TODO: Max session
	sessionId := 12345678

	// Generate an RNG message to send to NS
	rngMsg := RNGMessage{
		SwitchboardServerAddress: fmt.Sprintf("%s:%d", cf.SwitchboardServer.ServerAddr, cf.SwitchboardServer.ServerPort),
		SessionID:                uint32(sessionId),
		CallerEmail:              user.Email,
		CallerDisplayName:        user.DisplayName,
		CalleeEmail:              callee.Email,
	}

	jsonRngMsg, err := json.Marshal(rngMsg)
	if err != nil {
		return err
	}

	// Publish the message to inform NS
	if err := rdb.Publish(context.TODO(), cf.Redis.PubSubChannel, jsonRngMsg).Err(); err != nil {
		return err
	}

	res := fmt.Sprintf("CAL %d RINGING %d\r\n", tid, sessionId)
	c.Send(res)

	return nil
}
