package sessions

import (
	"msnserver/pkg/clients"
	"sync"
)

type SwitchboardSessions struct {
	sessions       map[uint32][]*clients.Client
	clientsSession map[string]uint32
	maxID          uint32
	m              *sync.Mutex
}

func NewSwitchboardSessions() *SwitchboardSessions {
	return &SwitchboardSessions{
		sessions:       map[uint32][]*clients.Client{},
		clientsSession: map[string]uint32{},
		maxID:          0,
		m:              &sync.Mutex{},
	}
}

func (sbs *SwitchboardSessions) CreateSession(c *clients.Client) uint32 {
	sbs.m.Lock()
	defer sbs.m.Unlock()

	sbs.maxID++
	sbs.sessions[sbs.maxID] = []*clients.Client{c}
	sbs.clientsSession[c.Id] = sbs.maxID

	return sbs.maxID
}

func (sbs *SwitchboardSessions) GetSessionID(c *clients.Client) uint32 {
	sbs.m.Lock()
	defer sbs.m.Unlock()

	return sbs.clientsSession[c.Id]
}
