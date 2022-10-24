package session

import (
	"context"
	"github.com/lonng/nano/internal/log"
)

const (
	ContextSessionKey = "ctx_session"
)

func CtxSetSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, ContextSessionKey, s)
}

func CtxGetSession(ctx context.Context) Session {
	s := ctx.Value(ContextSessionKey)
	if s == nil {
		log.Println("ctx doesn't contain a session, are you calling GetSessionFromCtx from inside a remote?")
		return nil
	}
	return s.(Session)
}
