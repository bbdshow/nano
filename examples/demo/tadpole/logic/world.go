package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/examples/demo/tadpole/logic/protocol"
	"github.com/lonng/nano/session"
)

// World contains all tadpoles
type World struct {
	component.Base
	*nano.Group
}

// NewWorld returns a world instance
func NewWorld() *World {
	return &World{
		Group: nano.NewGroup(uuid.New().String()),
	}
}

// Init initialize world component
func (w *World) Init() {
	session.Lifetime.OnClosed(func(s session.Session) {
		w.Leave(s)
		w.Broadcast("leave", &protocol.LeaveWorldResponse{ID: s.ID()})
		log.Println(fmt.Sprintf("session count: %d", w.Count()))
	})
}

// Enter was called when new guest enter
func (w *World) Enter(ctx context.Context, msg []byte) error {
	s := session.CtxGetSession(ctx)
	w.Add(s)
	log.Println(fmt.Sprintf("session count: %d", w.Count()))
	return s.Response(ctx, &protocol.EnterWorldResponse{ID: s.ID()})
}

// Update refresh tadpole's position
func (w *World) Update(ctx context.Context, msg []byte) error {
	return w.Broadcast("update", msg)
}

// Message handler was used to communicate with each other
func (w *World) Message(ctx context.Context, msg *protocol.WorldMessage) error {
	s := session.CtxGetSession(ctx)
	msg.ID = s.ID()
	return w.Broadcast("message", msg)
}
