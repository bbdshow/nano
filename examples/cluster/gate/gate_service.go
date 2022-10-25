package gate

import (
	"context"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/examples/cluster/protocol"
	"github.com/lonng/nano/session"
	"github.com/pingcap/errors"
)

type BindService struct {
	component.Base
	nextGateUid int64
}

func newBindService() *BindService {
	return &BindService{}
}

type (
	LoginRequest struct {
		Nickname string `json:"nickname"`
	}
	LoginResponse struct {
		Code int `json:"code"`
	}
)

func (bs *BindService) Login(ctx context.Context, msg *LoginRequest) error {
	s := session.CtxGetSession(ctx)
	bs.nextGateUid++
	uid := bs.nextGateUid
	request := &protocol.NewUserRequest{
		Nickname: msg.Nickname,
		GateUid:  uid,
	}
	if err := s.RPC(ctx, "TopicService.NewUser", request); err != nil {
		return errors.Trace(err)
	}
	return s.Response(ctx, &LoginResponse{})
}

func (bs *BindService) BindChatServer(s *session.Session, msg []byte) error {
	return errors.Errorf("not implement")
}
