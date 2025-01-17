package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/lonng/nano/cluster/clusterpb"
	"github.com/lonng/nano/timer"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/internal/env"
	"github.com/lonng/nano/internal/log"
	"github.com/lonng/nano/internal/message"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/session"
	"google.golang.org/grpc"
)

// Options contains some configurations for current node
type Options struct {
	Pipeline            pipeline.Pipeline
	IsMaster            bool
	AdvertiseAddr       string
	RetryInterval       time.Duration
	ClientAddr          string
	Components          *component.Components
	Label               string
	IsWebsocket         bool
	TSLCertificate      string
	TSLKey              string
	ConcurrentScheduler int // 并发调度器
	OnUnregister        func(serverAddr string)
}

//// OnUnregister 注销执行方法
//type OnUnregister func(serverAddr string)

// Node represents a node in nano cluster, which will contains a group of services.
// All services will register to cluster and messages will be forwarded to the node
// which provides respective service
type Node struct {
	Options            // current node options
	ServiceAddr string // current server service address (RPC)

	cluster   *cluster
	handler   *Handler
	server    *grpc.Server
	rpcClient *rpcClient

	mu       sync.RWMutex
	sessions map[int64]session.Session

	once          sync.Once
	heartbeatExit chan struct{}
}

func (n *Node) Startup() error {
	if n.ServiceAddr == "" {
		return errors.New("service address cannot be empty in master node")
	}
	n.sessions = map[int64]session.Session{}
	n.cluster = newCluster(n)
	n.handler = NewHandler(n, n.Pipeline)
	components := n.Components.List()
	for _, c := range components {
		err := n.handler.register(c.Comp, c.Opts)
		if err != nil {
			return err
		}
	}

	cache()
	if err := n.initNode(); err != nil {
		return err
	}

	// Initialize all components
	for _, c := range components {
		c.Comp.Init()
	}
	for _, c := range components {
		c.Comp.AfterInit()
	}

	if n.ClientAddr != "" {
		go func() {
			if n.IsWebsocket {
				if len(n.TSLCertificate) != 0 {
					n.listenAndServeWSTLS()
				} else {
					n.listenAndServeWS()
				}
			} else {
				n.listenAndServe()
			}
		}()
	}

	return nil
}

func (n *Node) Handler() *Handler {
	return n.handler
}

func (n *Node) initNode() error {
	// Current node is not master server and does not contains master
	// address, so running in singleton mode
	if !n.IsMaster && n.AdvertiseAddr == "" {
		return nil
	}

	listener, err := net.Listen("tcp", n.ServiceAddr)
	if err != nil {
		return err
	}

	// Initialize the gRPC server and register service
	n.server = grpc.NewServer()
	n.rpcClient = newRPCClient()
	clusterpb.RegisterMemberServer(n.server, n)

	go func() {
		err := n.server.Serve(listener)
		if err != nil {
			log.Fatalf("Start current node failed: %v", err)
		}
	}()

	if n.IsMaster {
		clusterpb.RegisterMasterServer(n.server, n.cluster)
		member := &Member{
			isMaster: true,
			memberInfo: &clusterpb.MemberInfo{
				Label:       n.Label,
				ServiceAddr: n.ServiceAddr,
				Services:    n.handler.LocalService(),
			},
		}
		n.cluster.members = append(n.cluster.members, member)
		n.cluster.setRpcClient(n.rpcClient)
	} else {
		pool, err := n.rpcClient.getConnPool(n.AdvertiseAddr)
		if err != nil {
			return err
		}
		client := clusterpb.NewMasterClient(pool.Get())
		request := &clusterpb.RegisterRequest{
			MemberInfo: &clusterpb.MemberInfo{
				Label:       n.Label,
				ServiceAddr: n.ServiceAddr,
				Services:    n.handler.LocalService(),
			},
		}
		for {
			resp, err := client.Register(context.Background(), request)
			if err == nil {
				n.handler.initRemoteService(resp.Members)
				n.cluster.initMembers(resp.Members)
				break
			}
			log.Println("Register current node to cluster failed", err, "and will retry in", n.RetryInterval.String())
			time.Sleep(n.RetryInterval)
		}
	}

	// 集群模式，开启心跳
	n.sendHeartbeat()

	return nil
}

func (n *Node) sendHeartbeat() {
	if n.heartbeatExit == nil {
		n.heartbeatExit = make(chan struct{})
	}
	n.once.Do(func() {
		go func() {
			ticker := time.NewTicker(env.Heartbeat)
			for {
				select {
				case <-ticker.C:
					n.heartbeat()
				case <-n.heartbeatExit:
					log.Println("member heartbeat exit")
					ticker.Stop()
					return
				}
			}
		}()
	})
}

// Shutdown all components registered by application, that
// call by reverse order against register
func (n *Node) Shutdown() {
	// reverse call `BeforeShutdown` hooks
	components := n.Components.List()
	length := len(components)
	for i := length - 1; i >= 0; i-- {
		components[i].Comp.BeforeShutdown()
	}

	// reverse call `Shutdown` hooks
	for i := length - 1; i >= 0; i-- {
		components[i].Comp.Shutdown()
	}
	if n.heartbeatExit != nil {
		close(n.heartbeatExit)
	}
	if !n.IsMaster && n.AdvertiseAddr != "" {
		pool, err := n.rpcClient.getConnPool(n.AdvertiseAddr)
		if err != nil {
			log.Println("Retrieve master address error", err)
			goto EXIT
		}
		client := clusterpb.NewMasterClient(pool.Get())
		request := &clusterpb.UnregisterRequest{
			ServiceAddr: n.ServiceAddr,
		}
		_, err = client.Unregister(context.Background(), request)
		if err != nil {
			log.Println("Unregister current node failed", err)
			goto EXIT
		}
	}

EXIT:
	if n.server != nil {
		n.server.GracefulStop()
	}
}

// Enable current server accept connection
func (n *Node) listenAndServe() {
	listener, err := net.Listen("tcp", n.ClientAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		go n.handler.handle(conn)
	}
}

func (n *Node) listenAndServeWS() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     env.CheckOrigin,
	}

	http.HandleFunc("/"+strings.TrimPrefix(env.WSPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error()))
			return
		}

		n.handler.handleWS(conn)
	})

	if err := http.ListenAndServe(n.ClientAddr, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func (n *Node) listenAndServeWSTLS() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     env.CheckOrigin,
	}

	http.HandleFunc("/"+strings.TrimPrefix(env.WSPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(fmt.Sprintf("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error()))
			return
		}

		n.handler.handleWS(conn)
	})

	if err := http.ListenAndServeTLS(n.ClientAddr, n.TSLCertificate, n.TSLKey, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func (n *Node) storeSession(s session.Session) {
	n.mu.Lock()
	n.sessions[s.ID()] = s
	n.mu.Unlock()
}

func (n *Node) findSession(sid int64) session.Session {
	n.mu.RLock()
	s := n.sessions[sid]
	n.mu.RUnlock()
	return s
}

func (n *Node) findOrCreateSession(sid int64, gateAddr string) (session.Session, error) {
	n.mu.RLock()
	s, found := n.sessions[sid]
	n.mu.RUnlock()
	if !found {
		//log.Println(sid, gateAddr, n.ServiceAddr)
		conns, err := n.rpcClient.getConnPool(gateAddr)
		if err != nil {
			return nil, err
		}
		ac := &acceptor{
			lastSessionId: sid,
			gateClient:    clusterpb.NewMemberClient(conns.Get()),
			rpcHandler:    n.handler.remoteProcess,
			gateAddr:      gateAddr,
		}
		s = session.NewSession(ac)

		ac.session = s

		n.mu.Lock()
		n.sessions[sid] = s
		n.mu.Unlock()
	}
	return s, nil
}

func (n *Node) HandleRequest(ctx context.Context, req *clusterpb.RequestMessage) (*clusterpb.MemberHandleResponse, error) {
	found := n.handler.IsLocalHandler(req.Route)
	if !found {
		return nil, fmt.Errorf("service not found in current node: %v", req.Route)
	}
	s, err := n.findOrCreateSession(req.SessionId, req.GateAddr)
	if err != nil {
		return nil, err
	}
	msg := &message.Message{
		Type:  message.Request,
		ID:    req.Id,
		Route: req.Route,
		Data:  req.Data,
	}

	ctx = session.CtxSetSession(context.Background(), s)

	uMsg := unhandledMessage{
		ctx: ctx,
		msg: msg,
		mid: req.Id,
	}
	n.handler.sendUnhandledMessage(uMsg, false)
	return &clusterpb.MemberHandleResponse{}, nil
}

func (n *Node) HandleNotify(ctx context.Context, req *clusterpb.NotifyMessage) (*clusterpb.MemberHandleResponse, error) {
	//log.Println("HandleNotify", req.Route, req.SessionId, req.GateAddr, n.ServiceAddr, string(req.Data))
	found := n.handler.IsLocalHandler(req.Route)
	if !found {
		return nil, fmt.Errorf("service not found in current node: %v", req.Route)
	}
	s, err := n.findOrCreateSession(req.SessionId, req.GateAddr)
	if err != nil {
		return nil, err
	}

	msg := &message.Message{
		Type:  message.Notify,
		Route: req.Route,
		Data:  req.Data,
	}
	ctx = session.CtxSetSession(context.Background(), s)

	uMsg := unhandledMessage{
		ctx: ctx,
		msg: msg,
		mid: 0,
	}
	n.handler.sendUnhandledMessage(uMsg, false)

	return &clusterpb.MemberHandleResponse{}, nil
}

func (n *Node) HandlePush(ctx context.Context, req *clusterpb.PushMessage) (*clusterpb.MemberHandleResponse, error) {
	s := n.findSession(req.SessionId)
	if s == nil {
		return &clusterpb.MemberHandleResponse{}, fmt.Errorf("HandlePush session not found: %v", req.SessionId)
	}

	return &clusterpb.MemberHandleResponse{}, s.Push(req.Route, req.Data)
}

func (n *Node) HandleResponse(ctx context.Context, req *clusterpb.ResponseMessage) (*clusterpb.MemberHandleResponse, error) {
	s := n.findSession(req.SessionId)
	if s == nil {
		return &clusterpb.MemberHandleResponse{}, fmt.Errorf("HandleResponse session not found: %v", req.SessionId)
	}
	return &clusterpb.MemberHandleResponse{}, s.ResponseMID(ctx, req.Id, req.Data)
}

func (n *Node) NewMember(_ context.Context, req *clusterpb.NewMemberRequest) (*clusterpb.NewMemberResponse, error) {
	n.handler.addRemoteService(req.MemberInfo)
	n.cluster.addMember(req.MemberInfo)
	return &clusterpb.NewMemberResponse{}, nil
}

func (n *Node) DelMember(_ context.Context, req *clusterpb.DelMemberRequest) (*clusterpb.DelMemberResponse, error) {
	log.Println("node DelMember", req.String())
	n.handler.delMember(req.ServiceAddr)
	n.cluster.delMember(req.ServiceAddr)
	return &clusterpb.DelMemberResponse{}, nil
}

// SessionClosed implements the MemberServer interface
func (n *Node) SessionClosed(_ context.Context, req *clusterpb.SessionClosedRequest) (*clusterpb.SessionClosedResponse, error) {
	n.mu.Lock()
	s, found := n.sessions[req.SessionId]
	delete(n.sessions, req.SessionId)
	n.mu.Unlock()
	if found {
		timer.NewAfterTimer(100*time.Millisecond, func() {
			session.Lifetime.Close(s)
		})
	}
	return &clusterpb.SessionClosedResponse{}, nil
}

// CloseSession implements the MemberServer interface
func (n *Node) CloseSession(_ context.Context, req *clusterpb.CloseSessionRequest) (*clusterpb.CloseSessionResponse, error) {
	n.mu.Lock()
	s, found := n.sessions[req.SessionId]
	delete(n.sessions, req.SessionId)
	n.mu.Unlock()
	if found {
		s.Close()
	}
	return &clusterpb.CloseSessionResponse{}, nil
}

// 定时向master发送注册信息
func (n *Node) heartbeat() {
	if n.AdvertiseAddr == "" || n.IsMaster {
		return
	}
	pool, err := n.rpcClient.getConnPool(n.AdvertiseAddr)
	if err != nil {
		log.Println("rpcClient master conn", err)
		return
	}
	masterCli := clusterpb.NewMasterClient(pool.Get())
	if _, err := masterCli.Register(context.Background(), &clusterpb.RegisterRequest{
		MemberInfo: &clusterpb.MemberInfo{
			Label:       n.Label,
			ServiceAddr: n.ServiceAddr,
			Services:    n.handler.LocalService(),
		},
		IsHeartbeat: true,
	}); err != nil {
		log.Println("heartbeat Register", err)
	}
}
