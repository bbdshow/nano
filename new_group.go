package nano

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lonng/nano/internal/env"
	"github.com/lonng/nano/internal/log"
	"github.com/lonng/nano/internal/message"
	"github.com/lonng/nano/session"
)

const (
	groupStatusWorking = 0
	groupStatusClosed  = 1
)

// SessionFilter represents a filter which was used to filter session when Multicast,
// the session will receive the message while filter returns true.
type SessionFilter func(session.Session) bool

// Group represents a session group which used to manage a number of
// sessions, data send to the group will send to all session in it.
type Group struct {
	mu       sync.RWMutex
	status   int32                      // channel current status
	name     string                     // channel name
	sessions map[string]session.Session // session uid map to session instance
}

// NewGroup returns a new group instance
func NewGroup(n string) *Group {
	return &Group{
		status:   groupStatusWorking,
		name:     n,
		sessions: make(map[string]session.Session),
	}
}

// Member returns specified UID's session
func (c *Group) Member(uid string) (session.Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, s := range c.sessions {
		if s.UID() == uid {
			return s, nil
		}
	}

	return nil, ErrMemberNotFound
}

// Members returns all member's UID in current group
func (c *Group) Members() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var members []string
	for _, s := range c.sessions {
		members = append(members, s.UID())
	}

	return members
}

// Multicast  push  the message to the filtered clients
func (c *Group) Multicast(route string, v interface{}, filter SessionFilter) error {
	if c.isClosed() {
		return ErrClosedGroup
	}

	data, err := message.Serialize(v)
	if err != nil {
		return err
	}

	if env.Debug {
		log.Println(fmt.Sprintf("Multicast %s, Data=%+v", route, v))
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, s := range c.sessions {
		if !filter(s) {
			continue
		}
		if err = s.Push(route, data); err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}

// Broadcast push  the message(s) to  all members
func (c *Group) Broadcast(route string, v interface{}) error {
	if c.isClosed() {
		return ErrClosedGroup
	}

	data, err := message.Serialize(v)
	if err != nil {
		return err
	}

	if env.Debug {
		log.Println(fmt.Sprintf("Broadcast %s, Data=%+v", route, v))
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, s := range c.sessions {
		if err = s.Push(route, data); err != nil {
			log.Println(fmt.Sprintf("Session push message error, ID=%d, UID=%d, Error=%s", s.ID(), s.UID(), err.Error()))
		}
	}

	return err
}

// Contains check whether a UID is contained in current group or not
func (c *Group) Contains(uid string) bool {
	_, err := c.Member(uid)
	return err == nil
}

// Add session to group
func (c *Group) Add(session session.Session) error {
	if c.isClosed() {
		return ErrClosedGroup
	}

	if env.Debug {
		log.Println(fmt.Sprintf("Add session to group %s, ID=%d, UID=%s", c.name, session.ID(), session.UID()))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	uid := session.UID()
	if uid == "" {
		return ErrUidEmpty
	}
	_, ok := c.sessions[uid]
	if ok {
		return ErrSessionDuplication
	}

	c.sessions[uid] = session
	return nil
}

// Leave remove specified UID related session from group
func (c *Group) Leave(s session.Session) error {
	if c.isClosed() {
		return ErrClosedGroup
	}

	if env.Debug {
		log.Println(fmt.Sprintf("Remove session from group %s, UID=%s", c.name, s.UID()))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.sessions, s.UID())
	return nil
}

// LeaveAll clear all sessions in the group
func (c *Group) LeaveAll() error {
	if c.isClosed() {
		return ErrClosedGroup
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.sessions = make(map[string]session.Session)
	return nil
}

// Count get current member amount in the group
func (c *Group) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.sessions)
}

func (c *Group) isClosed() bool {
	if atomic.LoadInt32(&c.status) == groupStatusClosed {
		return true
	}
	return false
}

// Close destroy group, which will release all resource in the group
func (c *Group) Close() error {
	if c.isClosed() {
		return ErrCloseClosedGroup
	}

	atomic.StoreInt32(&c.status, groupStatusClosed)

	// release all reference
	c.sessions = make(map[string]session.Session)
	return nil
}
