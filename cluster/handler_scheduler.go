package cluster

import (
	"context"
	"fmt"
	"github.com/lonng/nano/cluster/clusterpb"
	"github.com/lonng/nano/internal/env"
	"github.com/lonng/nano/internal/log"
	"github.com/lonng/nano/internal/message"
	"github.com/lonng/nano/session"
	"math/rand"
	"reflect"
	"runtime/debug"
	"strings"
)

func (h *Handler) Scheduler(n int) {
	for {
		select {
		case localMsg := <-h.chanLocalMsg:
			h.localProcess(localMsg.ctx, localMsg.mid, localMsg.msg)
		case remoteMsg := <-h.chanRemoteMsg:
			h.remoteProcess(remoteMsg.ctx, remoteMsg.msg)
		}
	}
}

func (h *Handler) localProcess(ctx context.Context, lastMid uint64, msg *message.Message) {
	sess := session.CtxGetSession(ctx)

	handler := h.localHandlers[msg.Route]

	if pipe := h.pipeline; pipe != nil {
		err := pipe.Inbound().Process(sess, msg)
		if err != nil {
			log.Println("Pipeline process failed: " + err.Error())
			return
		}
	}

	var payload = msg.Data
	var data interface{}
	if handler.IsRawArg {
		data = payload
	} else {
		data = reflect.New(handler.Type.Elem()).Interface()
		err := env.Serializer.Unmarshal(payload, data)
		if err != nil {
			log.Println(fmt.Sprintf("Deserialize to %T failed: %+v (%v)", data, err, payload))
			return
		}
	}

	if env.Debug {
		log.Println(fmt.Sprintf("UID=%s, Message={%s}, Data=%+v", sess.UID(), msg.String(), data))
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle message panic: %+v\n%s", err, debug.Stack()))
		}
	}()

	index := strings.LastIndex(msg.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("nano/handler: invalid route %s", msg.Route))
		return
	}

	args := []reflect.Value{handler.Receiver, reflect.ValueOf(ctx), reflect.ValueOf(data)}
	switch v := sess.NetworkEntity().(type) {
	case *agent:
		v.lastMid = lastMid
	case *acceptor:
		v.lastMid = lastMid
	}
	result := handler.Method.Func.Call(args)
	if len(result) > 0 {
		if err := result[0].Interface(); err != nil {
			log.Println(fmt.Sprintf("Service %s error: %+v", msg.Route, err))
		}
	}

	// A message can be dispatch to global thread or a user customized thread
	//service := msg.Route[:index]
	//log.Println("service", service, h.localServices[service])
	//
	//if s, found := h.localServices[service]; found && s.SchedName != "" {
	//	log.Println("s.SchedName", s.SchedName)
	//	sched := sess.Value(s.SchedName)
	//	if sched == nil {
	//		log.Println(fmt.Sprintf("nanl/handler: cannot found `schedular.LocalScheduler` by %s", s.SchedName))
	//		return
	//	}
	//
	//	local, ok := sched.(scheduler.LocalScheduler)
	//	if !ok {
	//		log.Println(fmt.Sprintf("nanl/handler: Type %T does not implement the `schedular.LocalScheduler` interface",
	//			sched))
	//		return
	//	}
	//	local.Schedule(func() {
	//
	//	})
	//} else {
	//
	//}
}

func (h *Handler) remoteProcess(ctx context.Context, msg *message.Message) {
	log.Println(msg.Route, msg.String(), h.currentNode.ServiceAddr)
	sess := session.CtxGetSession(ctx)
	index := strings.LastIndex(msg.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("nano/handler: invalid route %s", msg.Route))
		return
	}

	service := msg.Route[:index]
	members := h.findRemoteMembers(service)
	if len(members) == 0 {
		log.Println(fmt.Sprintf("nano/handler: %s not found(forgot registered?)", msg.Route))
		return
	}

	// Select a remote service address
	// 1. Use the service address directly if the router contains binding item
	// 2. Select a remote service address randomly and bind to router
	var remoteAddr string
	if addr, found := sess.Router().Find(service); found {
		remoteAddr = addr
	} else {
		remoteAddr = members[rand.Intn(len(members))].ServiceAddr
		sess.Router().Bind(service, remoteAddr)
	}
	pool, err := h.currentNode.rpcClient.getConnPool(remoteAddr)
	if err != nil {
		log.Println(err)
		return
	}

	var data = msg.Data
	// Retrieve gate address and session id
	gateAddr := h.currentNode.ServiceAddr
	sid := sess.ID()
	switch v := sess.NetworkEntity().(type) {
	case *acceptor:
		gateAddr = v.gateAddr
		sid = v.lastSessionId
	}
	log.Println("remoteAddr", remoteAddr, gateAddr, "ctx", &ctx)

	client := clusterpb.NewMemberClient(pool.Get())
	switch msg.Type {
	case message.Request:
		request := &clusterpb.RequestMessage{
			GateAddr:  gateAddr,
			SessionId: sid,
			Id:        msg.ID,
			Route:     msg.Route,
			Data:      data,
		}
		deadline, ok := ctx.Deadline()
		log.Println("Request", remoteAddr, gateAddr, request, deadline.String(), ok, ctx.Err())
		_, err = client.HandleRequest(ctx, request)
	case message.Notify:
		request := &clusterpb.NotifyMessage{
			GateAddr:  gateAddr,
			SessionId: sid,
			Route:     msg.Route,
			Data:      data,
		}
		log.Println("remoteAddr", remoteAddr, gateAddr, request)
		_, err = client.HandleNotify(ctx, request)
	}
	if err != nil {
		log.Println(fmt.Sprintf("Process remote message (%d:%s) error: %+v", msg.ID, msg.Route, err))
	}
}
