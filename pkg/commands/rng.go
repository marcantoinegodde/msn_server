package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/utils"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RNGMessage struct {
	SwitchboardServerAddress string `json:"switchboard_server_address"`
	SessionID                uint32 `json:"session_id"`
	CallerEmail              string `json:"caller_email"`
	CallerDisplayName        string `json:"caller_display_name"`
	CalleeEmail              string `json:"callee_email"`
}

func HandleRNG(rdb *redis.Client, m *sync.Mutex, clients map[string]*clients.Client, rngMessage RNGMessage) error {
	random, err := utils.GenerateRandomString(25)
	if err != nil {
		return err
	}

	cki := cki{
		Cki:       random,
		SessionID: rngMessage.SessionID,
	}

	jsonCki, err := json.Marshal(cki)
	if err != nil {
		return err
	}

	if err := rdb.Set(context.TODO(), rngMessage.CalleeEmail, jsonCki, CKI_TIMEOUT).Err(); err != nil {
		return err
	}

	m.Lock()
	callee, ok := clients[rngMessage.CalleeEmail]
	if ok {
		res := fmt.Sprintf("RNG %d %s %s %s %s %s\r\n", rngMessage.SessionID, rngMessage.SwitchboardServerAddress,
			SB_SECURITY_PACKAGE, cki.Cki, rngMessage.CallerEmail, rngMessage.CallerDisplayName)
		callee.Send(res)
	}
	m.Unlock()

	return nil
}
