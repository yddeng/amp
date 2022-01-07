package exec

import (
	"bytes"
	"github.com/yddeng/dnet/drpc"
	"initial-server/logger"
	"initial-server/protocol"
	"os/exec"
	"syscall"
)

type Application struct {
	AppID    int32 `json:"app_id"`
	Pid      int   `json:"pid"`
	CreateAt int64 `json:"create_at"`
}

func (app *Application) isAlive() bool {
	if err := syscall.Kill(app.Pid, 0); err == nil {
		return true
	}
	return false
}

/*
func (er *Executor) onAppExec(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.AppExecReq)
	logger.GetSugar().Infof("onAppExec %v\n", msg)

	ecmd := exec.Command(msg.GetName(), msg.GetArgs()...)
	ecmd.Dir = msg.GetPath()

	errBuff := bytes.Buffer{}
	ecmd.Stderr = &errBuff
	outBuff := bytes.Buffer{}
	ecmd.Stdout = &outBuff

	cmd := CommandWithCmd(ecmd)
	if err := cmd.Run(func(cmd *Cmd, err error) {
		er.Submit(func() {
			if err != nil { // exit or signal
				_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
			} else {
				_ = replier.Reply(&protocol.CmdExecResp{ErrStr: errBuff.String(), OutStr: outBuff.String()}, nil)
			}
		})
	}); err != nil {
		_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
	}
}
*/

func (er *Executor) onCmdExec(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.CmdExecReq)
	logger.GetSugar().Infof("onCmdExec %v", msg)

	ecmd := exec.Command(msg.GetName(), msg.GetArgs()...)
	ecmd.Dir = msg.GetDir()

	//errBuff := bytes.Buffer{}
	outBuff := bytes.Buffer{}
	ecmd.Stderr = &outBuff
	ecmd.Stdout = &outBuff

	cmd := CommandWithCmd(ecmd)
	if err := cmd.Run(int(msg.GetTimeout()), func(cmd *Cmd, err error) {
		er.Submit(func() {
			if err != nil {
				// exit or signal
				if cmd.ProcessState().Exited() {
					// 执行出错
					_ = replier.Reply(&protocol.CmdExecResp{OutStr: outBuff.String()}, nil)
				} else {
					// 超时
					_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
				}
			} else {
				_ = replier.Reply(&protocol.CmdExecResp{OutStr: outBuff.String()}, nil)
			}
		})
	}); err != nil {
		_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
	}
}

/*
func (er *Executor) onStart(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.StartReq)
	logger.GetSugar().Infof("onStart %v\n", msg)

	itemID := msg.GetItemID()
	info, ok := apps[itemID]
	if ok && info.isAlive() {
		replier.Reply(&protocol.StartResp{Code: "itemID is started"}, nil)
		return
	}

	shell := fmt.Sprintf("nohup %s deploy > /dev/null 2> %s/item_%d.log & echo $!", msg.GetArgs(), er.cfg.FilePath, itemID)
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

	er.addExecInfo(itemID, pid)
	logger.GetSugar().Infof("start ok,itemID %d pid %d", itemID, pid)
	replier.Reply(&protocol.StartResp{}, nil)
}

func (er *Executor) onSignal(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.SignalReq)
	logger.GetSugar().Infof("onSignal %v\n", msg)

	itemID := msg.GetItemID()
	p, ok := er.execInfos[itemID]
	if !ok {
		replier.Reply(&protocol.SignalResp{Code: "itemID not exist"}, nil)
		return
	}

	signal := syscall.Signal(msg.GetSignal())
	var err error
	switch signal {
	case syscall.SIGTERM:
		if err = syscall.Kill(p.Pid, syscall.SIGTERM); err == nil {
			er.delExecInfo(itemID)
		}
	case syscall.SIGKILL:
		if err = syscall.Kill(p.Pid, syscall.SIGKILL); err == nil {
			er.delExecInfo(itemID)
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

func (er *Executor) onItemStatus(replier *drpc.Replier, req interface{}) {
	resp := &protocol.ItemStatueResp{
		Items: make(map[int32]*protocol.ItemStatue, len(er.execInfos)),
	}

	for _, v := range er.execInfos {
		resp.Items[v.ItemID] = &protocol.ItemStatue{
			ItemID:    v.ItemID,
			Pid:       int32(v.Pid),
			Timestamp: v.Timestamp,
			IsAlive:   v.isAlive(),
		}
	}
	replier.Reply(resp, nil)
}

func (er *Executor) onPanicLog(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.PanicLogReq)
	logger.GetSugar().Infof("onPanicLog %v\n", msg)

	filename := fmt.Sprintf("%s/item_%d.log", er.cfg.FilePath, msg.GetItemID())
	if data, err := ioutil.ReadFile(filename); err != nil {
		logger.GetSugar().Error(err)
		replier.Reply(&protocol.PanicLogResp{Code: err.Error()}, nil)
	} else {
		replier.Reply(&protocol.PanicLogResp{Msg: string(data)}, nil)
	}
}
*/
