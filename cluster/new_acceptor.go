package cluster

import (
	"context"
	"github.com/lonng/nano/cluster/clusterpb"
	"github.com/lonng/nano/internal/message"
	"github.com/lonng/nano/mock"
	"github.com/lonng/nano/session"
	"net"
)

type acceptor struct {
	gateClient clusterpb.MemberClient
	session    session.Session
	lastMid    uint64
	rpcHandler rpcHandler
	gateAddr   string
}

// Push implements the session.NetworkEntity interface
func (a *acceptor) Push(route string, v interface{}) error {
	// TODO: buffer
	data, err := message.Serialize(v)
	if err != nil {
		return err
	}
	request := &clusterpb.PushMessage{
		SessionUid: a.session.UID(),
		Route:      route,
		Data:       data,
	}
	_, err = a.gateClient.HandlePush(context.Background(), request)
	return err
}

// RPC implements the session.NetworkEntity interface
func (a *acceptor) RPC(ctx context.Context, route string, v interface{}) error {
	// TODO: buffer
	data, err := message.Serialize(v)
	if err != nil {
		return err
	}
	msg := &message.Message{
		Type:  message.Notify,
		Route: route,
		Data:  data,
	}
	a.rpcHandler(ctx, msg, true)
	return nil
}

// LastMid implements the session.NetworkEntity interface
func (a *acceptor) LastMid() uint64 {
	return a.lastMid
}

// Response implements the session.NetworkEntity interface
func (a *acceptor) Response(ctx context.Context, v interface{}) error {
	return a.ResponseMid(ctx, a.lastMid, v)
}

// ResponseMid implements the session.NetworkEntity interface
func (a *acceptor) ResponseMid(ctx context.Context, mid uint64, v interface{}) error {
	// TODO: buffer
	data, err := message.Serialize(v)
	if err != nil {
		return err
	}
	request := &clusterpb.ResponseMessage{
		SessionUid: a.session.UID(),
		Id:         mid,
		Data:       data,
	}
	_, err = a.gateClient.HandleResponse(ctx, request)
	return err
}

// Close implements the session.NetworkEntity interface
func (a *acceptor) Close() error {
	// TODO: buffer
	request := &clusterpb.CloseSessionRequest{
		SessionUid: a.session.UID(),
	}
	_, err := a.gateClient.CloseSession(context.Background(), request)
	return err
}

// RemoteAddr implements the session.NetworkEntity interface
func (*acceptor) RemoteAddr() net.Addr {
	return mock.NetAddr{}
}