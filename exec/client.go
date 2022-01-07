package exec

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/utils/task"
	"initial-server/logger"
	"initial-server/protocol"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Client struct {
	cfg       *Config
	dialing   bool
	taskPool  *task.TaskPool
	session   dnet.Session
	rpcServer *drpc.Server
	rpcClient *drpc.Client

	execInfos    map[int32]*execInfo
	execFilename string
}

func (c *Client) SendRequest(req *drpc.Request) error {
	return c.session.Send(req)
}

func (c *Client) SendResponse(resp *drpc.Response) error {
	return c.session.Send(resp)
}

func (c *Client) Go(data proto.Message, callback func(interface{}, error)) error {
	return c.rpcClient.Go(c, proto.MessageName(data), data, time.Second*5, callback)
}

/*
func (this *Client) Call(data proto.Message) (interface{}, error) {
	return this.rpcClient.Call(this, proto.MessageName(data), data, time.Second*5)
}
*/
func (c *Client) dial() {
	if c.session != nil || c.dialing {
		return
	}

	c.dialing = true

	go func() {
		for {
			conn, err := dnet.DialTCP(c.cfg.Center, time.Second*5)
			if nil == err {
				c.onConnected(conn)
				return
			} else {
				log.Println(err)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (c *Client) onConnected(conn net.Conn) {
	c.taskPool.Submit(func() {
		c.dialing = false
		c.session = dnet.NewTCPSession(conn,
			dnet.WithCodec(new(protocol.Codec)),
			dnet.WithMessageCallback(func(session dnet.Session, data interface{}) {
				c.taskPool.Submit(func() {
					switch data.(type) {
					case *drpc.Request:
						c.rpcServer.OnRPCRequest(c, data.(*drpc.Request))
					case *drpc.Response:
						c.rpcClient.OnRPCResponse(data.(*drpc.Response))
					case *protocol.Message:
						c.dispatchMsg(session, data.(*protocol.Message))
					}
				})
			}),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				c.taskPool.Submit(func() {
					c.session = nil
					logger.GetSugar().Infof("session closed, reason: %s\n", reason)
					time.Sleep(time.Second)
					c.dial()
				})

			}))

		// login
		if err := c.Go(&protocol.LoginReq{
			Name: c.cfg.Name,
			Net:  c.cfg.Net,
			Inet: c.cfg.Inet,
		}, func(i interface{}, e error) {
			if e != nil {
				c.session.Close(e)
				return
			}
			resp := i.(*protocol.LoginResp)
			if resp.GetCode() != "" {
				e = errors.New(resp.GetCode())
				c.session.Close(e)
				panic(e)
			}
		}); err != nil {
			c.session.Close(err)
		}

	})

}

func (c *Client) dispatchMsg(session dnet.Session, msg *protocol.Message) {}

type Config struct {
	Name     string
	Inet     string
	Net      string
	Center   string
	FilePath string
}

func NewClient(cfg Config) *Client {
	c := new(Client)
	c.cfg = &cfg
	c.taskPool = task.NewTaskPool(1, 1024)
	c.rpcClient = drpc.NewClient()
	c.rpcServer = drpc.NewServer()

	os.MkdirAll(cfg.FilePath, os.ModePerm)
	c.execFilename = path.Join(cfg.FilePath, dataFile)
	c.execInfos = loadExecInfo(c.execFilename)

	c.rpcServer.Register(proto.MessageName(&protocol.StartReq{}), c.onStart)
	c.rpcServer.Register(proto.MessageName(&protocol.SignalReq{}), c.onSignal)
	c.rpcServer.Register(proto.MessageName(&protocol.ItemStatueReq{}), c.onItemStatus)
	c.rpcServer.Register(proto.MessageName(&protocol.PanicLogReq{}), c.onPanicLog)
	return c
}

func (c *Client) Start() {
	c.taskPool.Submit(c.dial)
}

func (c *Client) onStart(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.StartReq)
	logger.GetSugar().Infof("onStart %v\n", msg)

	itemID := msg.GetItemID()
	info, ok := c.execInfos[itemID]
	if ok && info.isAlive() {
		replier.Reply(&protocol.StartResp{Code: "itemID is started"}, nil)
		return
	}

	shell := fmt.Sprintf("nohup %s deploy > /dev/null 2> %s/item_%d.log & echo $!", msg.GetArgs(), c.cfg.FilePath, itemID)
	logger.GetSugar().Debug(itemID, shell)
	cmd := exec.Command("sh", "-c", shell)
	out, err := cmd.Output()
	if err != nil {
		replier.Reply(&protocol.StartResp{Code: err.Error()}, nil)
		logger.GetSugar().Errorf(err.Error())
		return
	}

	// 进程pid
	str := strings.Split(string(out), "\n")[0]
	pid, err := strconv.Atoi(str)
	if nil != err {
		logger.GetSugar().Error("strconv.Atoi pid error:", string(out), err)
		replier.Reply(&protocol.StartResp{Code: "parseInt pid error"}, nil)
		return
	}

	c.addExecInfo(itemID, pid)
	logger.GetSugar().Infof("start ok,itemID %d pid %d", itemID, pid)
	replier.Reply(&protocol.StartResp{}, nil)
}

func (c *Client) onSignal(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.SignalReq)
	logger.GetSugar().Infof("onSignal %v\n", msg)

	itemID := msg.GetItemID()
	p, ok := c.execInfos[itemID]
	if !ok {
		replier.Reply(&protocol.SignalResp{Code: "itemID not exist"}, nil)
		return
	}

	signal := syscall.Signal(msg.GetSignal())
	var err error
	switch signal {
	case syscall.SIGTERM:
		if err = syscall.Kill(p.Pid, syscall.SIGTERM); err == nil {
			c.delExecInfo(itemID)
		}
	case syscall.SIGKILL:
		if err = syscall.Kill(p.Pid, syscall.SIGKILL); err == nil {
			c.delExecInfo(itemID)
		}
	default:
		err = syscall.Kill(p.Pid, signal)
	}
	if err != nil {
		replier.Reply(&protocol.SignalResp{Code: err.Error()}, nil)
		logger.GetSugar().Errorf(err.Error())
		return
	}

	logger.GetSugar().Infof("signal ok, itemID %d", itemID)
	replier.Reply(&protocol.SignalResp{}, nil)
}

func (c *Client) onItemStatus(replier *drpc.Replier, req interface{}) {
	resp := &protocol.ItemStatueResp{
		Items: make(map[int32]*protocol.ItemStatue, len(c.execInfos)),
	}

	for _, v := range c.execInfos {
		resp.Items[v.ItemID] = &protocol.ItemStatue{
			ItemID:    v.ItemID,
			Pid:       int32(v.Pid),
			Timestamp: v.Timestamp,
			IsAlive:   v.isAlive(),
		}
	}
	replier.Reply(resp, nil)
}

func (c *Client) onPanicLog(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.PanicLogReq)
	logger.GetSugar().Infof("onPanicLog %v\n", msg)

	filename := fmt.Sprintf("%s/item_%d.log", c.cfg.FilePath, msg.GetItemID())
	if data, err := ioutil.ReadFile(filename); err != nil {
		logger.GetSugar().Error(err)
		replier.Reply(&protocol.PanicLogResp{Code: err.Error()}, nil)
	} else {
		replier.Reply(&protocol.PanicLogResp{Msg: string(data)}, nil)
	}
}
