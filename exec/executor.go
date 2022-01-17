package exec

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/utils/task"
	"initial-server/logger"
	"initial-server/protocol"
	"math/rand"
	"net"
	"time"
)

type Config struct {
	Name     string `json:"name"`
	Inet     string `json:"inet"`
	Net      string `json:"net"`
	Center   string `json:"center"`
	DataPath string `json:"data_path"`
}

type Executor struct {
	cfg       *Config
	dialing   bool
	session   dnet.Session
	taskPool  *task.TaskPool
	rpcServer *drpc.Server
	rpcClient *drpc.Client
}

func (er *Executor) SendRequest(req *drpc.Request) error {
	if er.session == nil {
		return errors.New("session is nil")
	}
	return er.session.Send(req)
}

func (er *Executor) SendResponse(resp *drpc.Response) error {
	if er.session == nil {
		return errors.New("session is nil")
	}
	return er.session.Send(resp)
}

func (er *Executor) Go(data proto.Message, callback func(interface{}, error)) error {
	return er.rpcClient.Go(er, proto.MessageName(data), data, time.Second*5, callback)
}

func (er *Executor) Submit(fn interface{}, args ...interface{}) error {
	return er.taskPool.Submit(fn, args...)
}

func (er *Executor) dial() {
	if er.session != nil || er.dialing {
		return
	}

	er.dialing = true

	go func() {
		for {
			conn, err := dnet.DialTCP(er.cfg.Center, time.Second*5)
			if nil == err {
				er.onConnected(conn)
				return
			} else {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+500))
			}
		}
	}()
}

func (er *Executor) onConnected(conn net.Conn) {
	er.Submit(func() {
		er.dialing = false
		er.session = dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				er.Submit(func() {
					switch data.(type) {
					case *drpc.Request:
						er.rpcServer.OnRPCRequest(er, data.(*drpc.Request))
					case *drpc.Response:
						er.rpcClient.OnRPCResponse(data.(*drpc.Response))
					case *protocol.Message:
						er.dispatchMsg(session, data.(*protocol.Message))
					}
				})
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				er.Submit(func() {
					er.session.SetContext(nil)
					er.session = nil
					logger.GetSugar().Infof("session closed, reason: %s", reason)
					er.dial()
				})

			}))

		// login
		if err := er.Go(&protocol.LoginReq{
			Name: er.cfg.Name,
			Net:  er.cfg.Net,
			Inet: er.cfg.Inet,
		}, func(i interface{}, err error) {
			if err != nil {
				er.session.Close(err)
				panic(err)
			}
			resp := i.(*protocol.LoginResp)
			if resp.GetCode() != "" {
				err = errors.New(resp.GetCode())
				er.session.Close(err)
				panic(err)
			}
		}); err != nil {
			er.session.Close(err)
			panic(err)
		}

	})

}

func (er *Executor) dispatchMsg(session dnet.Session, msg *protocol.Message) {}

var er *Executor

func Start(cfg Config) (err error) {
	er = new(Executor)
	er.cfg = &cfg
	er.taskPool = task.NewTaskPool(1, 1024)
	er.rpcClient = drpc.NewClient()
	er.rpcServer = drpc.NewServer()

	er.rpcServer.Register(proto.MessageName(&protocol.CmdExecReq{}), er.onCmdExec)
	//er.rpcServer.Register(proto.MessageName(&protocol.StartReq{}), er.onStart)
	//er.rpcServer.Register(proto.MessageName(&protocol.SignalReq{}), er.onSignal)
	//er.rpcServer.Register(proto.MessageName(&protocol.ItemStatueReq{}), er.onItemStatus)
	//er.rpcServer.Register(proto.MessageName(&protocol.PanicLogReq{}), er.onPanicLog)

	er.Submit(er.dial)

	return nil
}
