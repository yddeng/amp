package exec

import (
	"amp/common"
	"amp/protocol"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/utils/task"
	"log"
	"math/rand"
	"net"
	"time"
)

type Config struct {
	Name     string `json:"name"`
	Inet     string `json:"inet"`
	Net      string `json:"net"`
	Center   string `json:"center"`
	Token    string `json:"token"`
	DataPath string `json:"data_path"`
}

type Executor struct {
	cfg       *Config
	dialing   bool
	session   dnet.Session
	taskPool  *task.TaskPool
	rpcServer *drpc.Server
	rpcClient *drpc.Client

	die            chan struct{}
	heartbeatTimer int64
}

func (er *Executor) SendMessage(msg proto.Message) error {
	if er.session == nil {
		return errors.New("session is nil")
	}
	return er.session.Send(protocol.NewMessage(msg))
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

func (er *Executor) closed() bool {
	select {
	case <-er.die:
		return true
	default:
		return false
	}
}

func (er *Executor) dial() {
	if er.session != nil || er.dialing || er.closed() {
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
		log.Printf("onConnected center %s", conn.RemoteAddr().String())
		er.dialing = false
		er.session = dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				er.Submit(func() {
					var err error
					switch data.(type) {
					case *drpc.Request:
						err = er.rpcServer.OnRPCRequest(er, data.(*drpc.Request))
					case *drpc.Response:
						err = er.rpcClient.OnRPCResponse(data.(*drpc.Response))
					case *protocol.Message:
						er.dispatchMsg(session, data.(*protocol.Message))
					}
					if err != nil {
						log.Println(err)
					}
				})
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				er.Submit(func() {
					er.session.SetContext(nil)
					er.session = nil
					log.Printf("session closed, reason: %s", reason)
					er.dial()
				})

			}))

		// login
		if err := er.Go(&protocol.LoginReq{
			Name:  er.cfg.Name,
			Net:   er.cfg.Net,
			Inet:  er.cfg.Inet,
			Token: er.cfg.Token,
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

func (er *Executor) tick() {
	timer := time.NewTimer(time.Second)
	heartbeatMsg := &protocol.Heartbeat{}
	for {
		select {
		case <-er.die:
			timer.Stop()
			return
		case now := <-timer.C:

			nodeStateMsg := packCollector()

			er.Submit(func() {
				_ = er.SendMessage(nodeStateMsg)
				if now.Unix() > er.heartbeatTimer {
					_ = er.SendMessage(heartbeatMsg)
					er.heartbeatTimer = now.Add(common.HeartbeatTimeout / 2).Unix()
				}
				timer.Reset(time.Second)
			})

		}
	}
}

func (er *Executor) dispatchMsg(session dnet.Session, msg *protocol.Message) {}

var er *Executor

func Start(cfg Config) (err error) {
	er = new(Executor)
	er.cfg = &cfg
	er.die = make(chan struct{})
	er.taskPool = task.NewTaskPool(1, 1024)
	er.rpcClient = drpc.NewClient()
	er.rpcServer = drpc.NewServer()

	er.rpcServer.Register(proto.MessageName(&protocol.CmdExecReq{}), er.onCmdExec)
	er.rpcServer.Register(proto.MessageName(&protocol.ProcessExecReq{}), er.onProcExec)
	er.rpcServer.Register(proto.MessageName(&protocol.ProcessSignalReq{}), er.onProcSignal)
	er.rpcServer.Register(proto.MessageName(&protocol.ProcessStateReq{}), er.onProcState)
	er.rpcServer.Register(proto.MessageName(&protocol.LogFileReq{}), er.onLogFile)

	loadProcess(cfg.DataPath)

	initCollector()

	er.Submit(er.dial)

	go er.tick()

	return nil
}

func Stop() {
	stopCh := make(chan struct{})
	er.Submit(func() {
		close(er.die)
		er.session.Close(fmt.Errorf("stop"))
		saveProcess()
		stopCh <- struct{}{}
	})

	//go func() {
	//	ticker := time.NewTicker(time.Millisecond * 50)
	//	for {
	//		<-ticker.C
	//		if er.taskPool.NumTask() == 0 {
	//			ticker.Stop()
	//			stopCh <- struct{}{}
	//		}
	//	}
	//}()
	<-stopCh
}
