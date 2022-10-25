package session

import (
	"context"
	"github.com/lonng/nano/internal/log"
	"net"
	"sync"
	"time"
)

/*
Fork Fix:
1. uid use string replace int64, user ensure uid cluster global unique
2. close session not use id as key, to impl cluster gate multi node。
*/

// Session represents a client session which could storage temp data during low-level
// keep connected, all data will be released when the low-level connection was broken.
// Session instance related to the client will be passed to Handler method as the first
type Session interface {
	NetworkEntity() NetworkEntity
	Router() *Router
	RPC(ctx context.Context, route string, v interface{}) error
	Push(route string, v interface{}) error
	Response(ctx context.Context, v interface{}) error
	ResponseMID(ctx context.Context, mid uint64, v interface{}) error

	ID() int64
	UID() string
	Bind(uid string) error

	LastMid() uint64

	RemoteAddr() net.Addr

	SetData(data map[string]interface{})
	GetData() map[string]interface{}

	Remove(key string)
	Set(key string, value interface{})
	HasKey(key string) bool

	Int(key string) int
	Int8(key string) int8
	Int16(key string) int16
	Int32(key string) int32
	Int64(key string) int64
	Uint(key string) uint
	Uint8(key string) uint8
	Uint16(key string) uint16
	Uint32(key string) uint32
	Uint64(key string) uint64
	Float32(key string) float32
	Float64(key string) float64
	String(key string) string
	Value(key string) interface{}
	Clear()

	Close()
}

// NewSession returns a new session instance
// a NetworkEntity is a low-level network instance
func NewSession(entity NetworkEntity) Session {
	s := &sessionImpl{
		id:       UniqueSessionId(),
		entity:   entity,
		data:     make(map[string]interface{}),
		lastTime: time.Now().Unix(),
		router:   newRouter(),
	}

	return s
}

// session impl
type sessionImpl struct {
	sync.RWMutex
	id       int64                  // session global unique id
	uid      string                 // binding user id
	lastTime int64                  // last heartbeat time
	entity   NetworkEntity          // low-level network entity
	data     map[string]interface{} // session data store
	router   *Router
}

// NetworkEntity returns the low-level network agent object
func (s *sessionImpl) NetworkEntity() NetworkEntity {
	return s.entity
}

// Router returns the service router
func (s *sessionImpl) Router() *Router {
	return s.router
}

// RPC sends message to remote server
func (s *sessionImpl) RPC(ctx context.Context, route string, v interface{}) error {
	return s.entity.RPC(ctx, route, v)
}

// Push message to client
func (s *sessionImpl) Push(route string, v interface{}) error {
	return s.entity.Push(route, v)
}

// Response message to client
func (s *sessionImpl) Response(ctx context.Context, v interface{}) error {
	return s.entity.Response(ctx, v)
}

// ResponseMID responses message to client, mid is
// request message ID
func (s *sessionImpl) ResponseMID(ctx context.Context, mid uint64, v interface{}) error {
	return s.entity.ResponseMid(ctx, mid, v)
}

// ID returns the session id
func (s *sessionImpl) ID() int64 {
	return s.id
}

// UID returns uid that bind to current session
func (s *sessionImpl) UID() string {
	return s.uid
}

// LastMid returns the last message id
func (s *sessionImpl) LastMid() uint64 {
	return s.entity.LastMid()
}

// Bind  UID to current session
func (s *sessionImpl) Bind(uid string) error {
	if uid == "" {
		return ErrIllegalUID
	}
	s.uid = uid
	return nil
}

// Close terminate current session, session related data will not be released,
// all related data should be Clear explicitly in Session closed callback
func (s *sessionImpl) Close() {
	if err := s.entity.Close(); err != nil {
		log.Println("entity close ", err)
	}
}

// RemoteAddr returns the remote network address.
func (s *sessionImpl) RemoteAddr() net.Addr {
	return s.entity.RemoteAddr()
}

func (s *sessionImpl) SetData(data map[string]interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data = data
}

func (s *sessionImpl) GetData() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

// Remove delete data associated with the key from session storage
func (s *sessionImpl) Remove(key string) {
	s.Lock()
	defer s.Unlock()

	delete(s.data, key)
}

// Set associates value with the key in session storage
func (s *sessionImpl) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	s.data[key] = value
}

// HasKey decides whether a key has associated value
func (s *sessionImpl) HasKey(key string) bool {
	s.RLock()
	defer s.RUnlock()

	_, has := s.data[key]
	return has
}

// Int returns the value associated with the key as a int.
func (s *sessionImpl) Int(key string) int {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

// Int8 returns the value associated with the key as a int8.
func (s *sessionImpl) Int8(key string) int8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

// Int16 returns the value associated with the key as a int16.
func (s *sessionImpl) Int16(key string) int16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

// Int32 returns the value associated with the key as a int32.
func (s *sessionImpl) Int32(key string) int32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

// Int64 returns the value associated with the key as a int64.
func (s *sessionImpl) Int64(key string) int64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

// Uint returns the value associated with the key as a uint.
func (s *sessionImpl) Uint(key string) uint {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

// Uint8 returns the value associated with the key as a uint8.
func (s *sessionImpl) Uint8(key string) uint8 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

// Uint16 returns the value associated with the key as a uint16.
func (s *sessionImpl) Uint16(key string) uint16 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

// Uint32 returns the value associated with the key as a uint32.
func (s *sessionImpl) Uint32(key string) uint32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

// Uint64 returns the value associated with the key as a uint64.
func (s *sessionImpl) Uint64(key string) uint64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

// Float32 returns the value associated with the key as a float32.
func (s *sessionImpl) Float32(key string) float32 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

// Float64 returns the value associated with the key as a float64.
func (s *sessionImpl) Float64(key string) float64 {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return 0
	}

	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

// String returns the value associated with the key as a string.
func (s *sessionImpl) String(key string) string {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.data[key]
	if !ok {
		return ""
	}

	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

// Value returns the value associated with the key as a interface{}.
func (s *sessionImpl) Value(key string) interface{} {
	s.RLock()
	defer s.RUnlock()

	return s.data[key]
}

// Clear releases all data related to current session
func (s *sessionImpl) Clear() {
	s.Lock()
	defer s.Unlock()

	s.uid = ""
	s.data = map[string]interface{}{}
}
