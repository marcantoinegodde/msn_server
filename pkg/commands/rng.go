package commands

import (
	"context"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/utils"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RNGMessage struct {
	SwitchboardServerAddress string `json:"switchboard_server_address"`
	SessionID                string `json:"session_id"`
	CallerEmail              string `json:"caller_email"`
	CallerDisplayName        string `json:"caller_display_name"`
	CalleeEmail              string `json:"callee_email"`
}

func HandleRNG(rdb *redis.Client, m *sync.Mutex, clients map[string]*clients.Client, rngMessage RNGMessage) error {
	cki := utils.GenerateRandomString(25)
	if err := rdb.Set(context.TODO(), rngMessage.CalleeEmail, cki, CKI_TIMEOUT).Err(); err != nil {
		return err
	}

	m.Lock()
	callee, ok := clients[rngMessage.CalleeEmail]
	if ok {
		res := fmt.Sprintf("RNG %s %s %s %s %s %s\r\n", rngMessage.SessionID, rngMessage.SwitchboardServerAddress,
			SB_SECURITY_PACKAGE, cki, rngMessage.CallerEmail, rngMessage.CallerDisplayName)
		callee.Send(res)
	}
	m.Unlock()

	return nil
}
