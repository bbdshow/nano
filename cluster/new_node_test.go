package cluster

import (
	"context"
	"fmt"
	"github.com/lonng/nano/benchmark/io"
	"github.com/lonng/nano/benchmark/testdata"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	. "github.com/pingcap/check"
	"strings"
	"testing"
)

type nodeSuite struct{}

var _ = Suite(&nodeSuite{})

type (
	MasterComponent struct{ component.Base }
	GateComponent   struct{ component.Base }
	GameComponent   struct{ component.Base }
)

func (c *MasterComponent) Test(ctx context.Context, _ []byte) error {
	sess := session.CtxGetSession(ctx)
	return sess.Push("test", &testdata.Pong{Content: "master server pong"})
}

func (c *GateComponent) Test(ctx context.Context, ping *testdata.Ping) error {
	sess := session.CtxGetSession(ctx)
	return sess.Push("test", &testdata.Pong{Content: "gate server pong"})
}

func (c *GateComponent) Test2(ctx context.Context, ping *testdata.Ping) error {
	sess := session.CtxGetSession(ctx)
	return sess.Response(ctx, &testdata.Pong{Content: "gate server pong2"})
}

func (c *GameComponent) Test(ctx context.Context, _ []byte) error {
	fmt.Println("GameComponent.Test")
	sess := session.CtxGetSession(ctx)
	//sess.Set("GameComponent", c.Schedule)
	fmt.Println(sess.ID())
	//time.Sleep(time.Minute)
	return sess.Push("test", &testdata.Pong{Content: "game server pong"})
}

func (c *GameComponent) Test2(ctx context.Context, ping *testdata.Ping) error {
	sess := session.CtxGetSession(ctx)
	//panic("=========================")
	fmt.Println("test2", sess.ID())
	return sess.Response(ctx, &testdata.Pong{Content: "game server pong2"})
}

func TestNode(t *testing.T) {
	TestingT(t)
}

func (s *nodeSuite) TestNodeStartup(c *C) {
	masterComps := &component.Components{}
	masterComps.Register(&MasterComponent{})
	masterNode := &Node{
		Options: Options{
			IsMaster:   true,
			Components: masterComps,
		},
		ServiceAddr: "127.0.0.1:4450",
	}
	err := masterNode.Startup()
	c.Assert(err, IsNil)
	masterHandler := masterNode.Handler()
	c.Assert(masterHandler.LocalService(), DeepEquals, []string{"MasterComponent"})

	member1Comps := &component.Components{}
	member1Comps.Register(&GateComponent{})
	memberNode1 := &Node{
		Options: Options{
			AdvertiseAddr: "127.0.0.1:4450",
			ClientAddr:    "127.0.0.1:14452",
			Components:    member1Comps,
		},
		ServiceAddr: "127.0.0.1:14451",
	}
	err = memberNode1.Startup()
	c.Assert(err, IsNil)
	member1Handler := memberNode1.Handler()
	c.Assert(masterHandler.LocalService(), DeepEquals, []string{"MasterComponent"})
	c.Assert(masterHandler.RemoteService(), DeepEquals, []string{"GateComponent"})
	c.Assert(member1Handler.LocalService(), DeepEquals, []string{"GateComponent"})
	c.Assert(member1Handler.RemoteService(), DeepEquals, []string{"MasterComponent"})

	member2Comps := &component.Components{}
	member2Comps.Register(&GameComponent{}, component.WithSchedulerName("Scheduler"))
	memberNode2 := &Node{
		Options: Options{
			AdvertiseAddr: "127.0.0.1:4450",
			Components:    member2Comps,
		},
		ServiceAddr: "127.0.0.1:24451",
	}
	err = memberNode2.Startup()
	c.Assert(err, IsNil)
	member2Handler := memberNode2.Handler()
	c.Assert(masterHandler.LocalService(), DeepEquals, []string{"MasterComponent"})
	c.Assert(masterHandler.RemoteService(), DeepEquals, []string{"GameComponent", "GateComponent"})
	c.Assert(member1Handler.LocalService(), DeepEquals, []string{"GateComponent"})
	c.Assert(member1Handler.RemoteService(), DeepEquals, []string{"GameComponent", "MasterComponent"})
	c.Assert(member2Handler.LocalService(), DeepEquals, []string{"GameComponent"})
	c.Assert(member2Handler.RemoteService(), DeepEquals, []string{"GateComponent", "MasterComponent"})

	// 启动调度器
	go masterNode.Handler().Scheduler(0)
	go memberNode1.Handler().Scheduler(0)
	go memberNode2.Handler().Scheduler(0)

	//fmt.Println("====")
	connector := io.NewConnector()

	chWait := make(chan struct{})
	connector.OnConnected(func() {
		chWait <- struct{}{}
	})

	// Connect to gate server
	if err := connector.Start("127.0.0.1:14452"); err != nil {
		c.Assert(err, IsNil)
	}
	<-chWait
	onResult := make(chan string)
	connector.On("test", func(data interface{}) {
		onResult <- string(data.([]byte))
	})
	err = connector.Notify("GateComponent.Test", &testdata.Ping{Content: "ping"})
	c.Assert(err, IsNil)
	c.Assert(strings.Contains(<-onResult, "gate server pong"), IsTrue)

	err = connector.Notify("GameComponent.Test", &testdata.Ping{Content: "ping"})
	c.Assert(err, IsNil)
	c.Assert(strings.Contains(<-onResult, "game server pong"), IsTrue)

	err = connector.Request("GateComponent.Test2", &testdata.Ping{Content: "ping"}, func(data interface{}) {
		onResult <- string(data.([]byte))
	})
	c.Assert(err, IsNil)
	c.Assert(strings.Contains(<-onResult, "gate server pong2"), IsTrue)

	err = connector.Request("GameComponent.Test2", &testdata.Ping{Content: "ping"}, func(data interface{}) {
		onResult <- string(data.([]byte))
	})
	c.Assert(err, IsNil)
	c.Assert(strings.Contains(<-onResult, "game server pong2"), IsTrue)

	err = connector.Notify("MasterComponent.Test", &testdata.Ping{Content: "ping"})
	c.Assert(err, IsNil)
	c.Assert(strings.Contains(<-onResult, "master server pong"), IsTrue)
}
