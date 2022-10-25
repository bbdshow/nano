package master

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lonng/nano/component"
	"github.com/lonng/nano/examples/cluster/protocol"
	"github.com/lonng/nano/session"
	"github.com/pingcap/errors"
)

type User struct {
	session  session.Session
	nickname string
	gateId   int64
	masterId int64
	balance  int64
	message  int
}

type TopicService struct {
	component.Base
	nextUid int64
	users   map[int64]*User
}

func newTopicService() *TopicService {
	return &TopicService{
		users: map[int64]*User{},
	}
}

type ExistsMembersResponse struct {
	Members string `json:"members"`
}

func (ts *TopicService) NewUser(ctx context.Context, msg *protocol.NewUserRequest) error {
	s := session.CtxGetSession(ctx)
	ts.nextUid++
	uid := ts.nextUid
	if err := s.Bind(fmt.Sprintf("%d", s)); err != nil {
		return errors.Trace(err)
	}

	var members []string
	for _, u := range ts.users {
		members = append(members, u.nickname)
	}
	err := s.Push("onMembers", &ExistsMembersResponse{Members: strings.Join(members, ",")})
	if err != nil {
		return errors.Trace(err)
	}

	user := &User{
		session:  s,
		nickname: msg.Nickname,
		gateId:   msg.GateUid,
		masterId: uid,
		balance:  1000,
	}
	ts.users[uid] = user

	chat := &protocol.JoinRoomRequest{
		Nickname:  msg.Nickname,
		GateUid:   msg.GateUid,
		MasterUid: uid,
	}
	return s.RPC(ctx, "RoomService.JoinRoom", chat)
}

type UserBalanceResponse struct {
	CurrentBalance int64 `json:"currentBalance"`
}

func (ts *TopicService) Stats(ctx context.Context, msg *protocol.MasterStats) error {
	s := session.CtxGetSession(ctx)
	// It's OK to use map without lock because of this service running in main thread
	user, found := ts.users[msg.Uid]
	if !found {
		return errors.Errorf("User not found: %v", msg.Uid)
	}
	user.message++
	user.balance--
	return s.Push("onBalance", &UserBalanceResponse{user.balance})
}

func (ts *TopicService) userDisconnected(s session.Session) {
	uid, _ := strconv.ParseInt(s.UID(), 10, 64)
	delete(ts.users, uid)
	log.Println("User session disconnected", s.UID())
}
