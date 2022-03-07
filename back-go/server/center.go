package server

import (
	"amp/back-go/common"
	"amp/back-go/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"log"
	"net"
	"time"
)

var (
	center *Center
)

type Center struct {
	token     string
	acceptor  dnet.Acceptor
	rpcServer *drpc.Server
	rpcClient *drpc.Client
}

func newCenter(address, token string) *Center {
	c := new(Center)
	c.acceptor = dnet.NewTCPAcceptor(address)
	c.rpcClient = drpc.NewClient()
	c.rpcServer = drpc.NewServer()
	c.rpcServer.Register(proto.MessageName(&protocol.LoginReq{}), c.onLogin)
	log.Printf("tcp server run %s.\n", address)
	return c
}

func (c *Center) Go(n *Node, data proto.Message, timeout time.Duration, callback func(interface{}, error)) error {
	return c.rpcClient.Go(n, proto.MessageName(data), data, timeout, callback)
}

func (c *Center) startListener() error {
	return c.acceptor.ServeFunc(func(conn net.Conn) {
		dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			dnet.WithTimeout(common.HeartbeatTimeout, 0),
			dnet.WithErrorCallback(func(session dnet.Session, err error) {
				log.Println(err)
				session.Close(err)
			}),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				taskQueue.Submit(func() {
					var err error
					switch data.(type) {
					case *drpc.Request:
						err = c.rpcServer.OnRPCRequest(&Node{session: session}, data.(*drpc.Request))
					case *drpc.Response:
						err = c.rpcClient.OnRPCResponse(data.(*drpc.Response))
					case *protocol.Message:
						c.dispatchMsg(session, data.(*protocol.Message))
					}
					if err != nil {
						log.Println(err)
					}
				})
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				taskQueue.Submit(func() {
					log.Printf("session closed, reason: %s\n", reason)
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

func (c *Center) dispatchMsg(session dnet.Session, msg *protocol.Message) {
	cmd := msg.GetCmd()
	switch cmd {
	case protocol.CmdHeartbeat:
		_ = session.Send(msg)
	case protocol.CmdNodeState:
		ctx := session.Context()
		if ctx != nil {
			node := ctx.(*Node)
			node.onNodeState(msg.GetData().(*protocol.NodeState))
		}
	default:

	}

}
