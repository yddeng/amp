package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/utils/task"
	"initial-sever/logger"
	"initial-sever/protocol"
	"log"
	"net"
	"time"
)

var (
	centerTaskQueue = task.NewTaskPool(1, 2048)
	center          *Center
)

type Center struct {
	acceptor  dnet.Acceptor
	rpcServer *drpc.Server
	rpcClient *drpc.Client
}

func newCenter(address string) *Center {
	c := new(Center)
	c.acceptor = dnet.NewTCPAcceptor(address)
	c.rpcClient = drpc.NewClient()
	c.rpcServer = drpc.NewServer()
	c.rpcServer.Register(proto.MessageName(&protocol.LoginReq{}), c.onLogin)
	log.Printf("tcp server run %s.\n", address)
	return c
}

func (c *Center) startListener() error {
	return c.acceptor.ServeFunc(func(conn net.Conn) {
		dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			dnet.WithTimeout(time.Second*5, 0),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				centerTaskQueue.Submit(func() {
					switch data.(type) {
					case *drpc.Request:
						c.rpcServer.OnRPCRequest(&Node{session: session}, data.(*drpc.Request))
					case *drpc.Response:
						c.rpcClient.OnRPCResponse(data.(*drpc.Response))
					case *protocol.Message:
						//c.dispatchMsg(session, data.(*protocol.Message))
					}
				})
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				centerTaskQueue.Submit(func() {
					logger.GetSugar().Infof("session closed, reason: %s\n", reason)
					ctx := session.Context()
					if ctx != nil {
						client := ctx.(*Node)
						client.session = nil
						session.SetContext(nil)
					}
				})
			}))
	})

}

func (this *Center) tick() {

}

func (this *Center) start() {
	go func() {
		if err := this.startListener(); err != nil {
			panic(err)
		}
	}()

	go func() {
		timer := time.NewTimer(time.Second)
		for {
			<-timer.C
			centerTaskQueue.Submit(func() {
				this.tick()
				timer.Reset(time.Second)
			})
		}
	}()
}
