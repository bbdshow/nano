package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/examples/demo/tadpole/logic/protocol"
	"github.com/lonng/nano/session"
)

// Manager component
type Manager struct {
	component.Base
}

// NewManager returns  a new manager instance
func NewManager() *Manager {
	return &Manager{}
}

// Login handler was used to guest login
func (m *Manager) Login(ctx context.Context, msg *protocol.JoyLoginRequest) error {
	s := session.CtxGetSession(ctx)
	log.Println(msg)
	id := s.ID()
	s.Bind(fmt.Sprintf("%d", id))
	return s.Response(ctx, protocol.LoginResponse{
		Status: protocol.LoginStatusSucc,
		ID:     id,
	})
}
