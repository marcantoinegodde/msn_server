package sessions

import (
	"errors"
	"msnserver/pkg/clients"
	"slices"
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

func (sbs *SwitchboardSessions) JoinSession(c *clients.Client, sessionID uint32) ([]*clients.Client, error) {
	sbs.m.Lock()
	defer sbs.m.Unlock()

	if _, ok := sbs.sessions[sessionID]; !ok {
		return nil, errors.New("session not found")
	}

	s := sbs.sessions[sessionID]
	sbs.sessions[sessionID] = append(sbs.sessions[sessionID], c)
	sbs.clientsSession[c.Id] = sessionID

	return s, nil
}

func (sbs *SwitchboardSessions) LeaveSession(c *clients.Client) ([]*clients.Client, error) {
	sbs.m.Lock()
	defer sbs.m.Unlock()

	sessionID, ok := sbs.clientsSession[c.Id]
	if !ok {
		return nil, errors.New("client not in session")
	}

	if len(sbs.sessions[sessionID]) == 1 {
		delete(sbs.sessions, sessionID)
		delete(sbs.clientsSession, c.Id)
		return nil, nil
	}

	for i, client := range sbs.sessions[sessionID] {
		if client == c {
			sbs.sessions[sessionID] = slices.Delete(sbs.sessions[sessionID], i, i+1)
			break
		}
	}

	delete(sbs.clientsSession, c.Id)

	return sbs.sessions[sessionID], nil
}
