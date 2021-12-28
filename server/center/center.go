package center

import (
	"errors"
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

var center *Center

type Node struct {
	Name    string
	session dnet.Session
}

func (c *Node) SendRequest(req *drpc.Request) error {
	return c.session.Send(req)
}

func (c *Node) SendResponse(resp *drpc.Response) error {
	return c.session.Send(resp)
}

type Center struct {
	acceptor  dnet.Acceptor
	taskPool  *task.TaskPool
	nodes     map[string]*Node
	rpcServer *drpc.Server
	rpcClient *drpc.Client
	itemMgr   *ItemMgr
}

func NewCenter(address string) *Center {
	c := new(Center)
	c.acceptor = dnet.NewTCPAcceptor(address)
	c.taskPool = task.NewTaskPool(1, 1024)
	c.rpcClient = drpc.NewClient()
	c.rpcServer = drpc.NewServer()
	c.nodes = map[string]*Node{}

	c.rpcServer.Register(proto.MessageName(&protocol.LoginReq{}), c.onLogin)
	log.Printf("tcp server run :%s.\n", address)
	return c
}

func (c *Center) startListener() error {
	return c.acceptor.ServeFunc(func(conn net.Conn) {
		dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			//dnet.WithTimeout(time.Second )
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				c.taskPool.Submit(func() {
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
				c.taskPool.Submit(func() {
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

func (this *Center) Start() {
	go func() {
		if err := this.startListener(); err != nil {
			panic(err)
		}
	}()

	go func() {
		timer := time.NewTimer(time.Second)
		for {
			<-timer.C
			this.taskPool.Submit(func() {
				this.tick()
				timer.Reset(time.Second)
			})
		}
	}()
}

func RunCenter(address string) {
	center = NewCenter(address)
	center.Start()
}

func (this *Center) onLogin(replier *drpc.Replier, req interface{}) {
	channel := replier.Channel
	msg := req.(*protocol.LoginReq)
	logger.GetSugar().Infof("onLogin %v\n", msg)

	name := msg.GetName()
	client := this.nodes[name]
	if client == nil {
		client = &Node{Name: name}
		this.nodes[name] = client
	}
	if client.session != nil {
		replier.Reply(&protocol.LoginResp{Code: "client already login. "}, nil)
		channel.(*Node).session.Close(errors.New("client already login. "))
		return
	}

	client.session = channel.(*Node).session
	replier.Reply(&protocol.LoginResp{}, nil)
}
