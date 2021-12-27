package center

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/utils/task"
	"initialtool/deploy/logger"
	"initialtool/deploy/protocol"
	"net"
	"os"
	"time"
)

type Config struct {
	Address  string
	FilePath string
}

type Client struct {
	Name    string
	session dnet.Session
}

func (c *Client) SendRequest(req *drpc.Request) error {
	return c.session.Send(req)
}

func (c *Client) SendResponse(resp *drpc.Response) error {
	return c.session.Send(resp)
}

type Center struct {
	cfg       *Config
	acceptor  dnet.Acceptor
	taskPool  *task.TaskPool
	clients   map[string]*Client
	rpcServer *drpc.Server
	rpcClient *drpc.Client
	itemMgr   *ItemMgr
}

func NewCenter(cfg Config) *Center {
	c := new(Center)
	c.cfg = &cfg
	c.acceptor = dnet.NewTCPAcceptor(cfg.Address)
	c.taskPool = task.NewTaskPool(1, 1024)
	c.rpcClient = drpc.NewClient()
	c.rpcServer = drpc.NewServer()
	c.clients = map[string]*Client{}
	_ = os.MkdirAll(cfg.FilePath, os.ModePerm)

	c.rpcServer.Register(proto.MessageName(&protocol.LoginReq{}), c.onLogin)
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
						c.rpcServer.OnRPCRequest(&Client{session: session}, data.(*drpc.Request))
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
						client := ctx.(*Client)
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

func (this *Center) onLogin(replier *drpc.Replier, req interface{}) {
	channel := replier.Channel
	msg := req.(*protocol.LoginReq)
	logger.GetSugar().Infof("onLogin %v\n", msg)

	name := msg.GetName()
	client := this.clients[name]
	if client == nil {
		client = &Client{Name: name}
		this.clients[name] = client
	}
	if client.session != nil {
		replier.Reply(&protocol.LoginResp{Code: "client already login. "}, nil)
		channel.(*Client).session.Close(errors.New("client already login. "))
		return
	}

	client.session = channel.(*Client).session
	replier.Reply(&protocol.LoginResp{}, nil)
}
