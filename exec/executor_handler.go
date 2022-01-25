package exec

import (
	"bytes"
	"github.com/yddeng/dnet/drpc"
	"initial-server/logger"
	"initial-server/protocol"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

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
					if err.Error() == "signal: killed" {
						// 超时 kill
						_ = replier.Reply(&protocol.CmdExecResp{Code: "执行超时，已终止"}, nil)
					} else {
						_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
					}
				}
			} else {
				_ = replier.Reply(&protocol.CmdExecResp{OutStr: outBuff.String()}, nil)
			}
		})
	}); err != nil {
		_ = replier.Reply(&protocol.CmdExecResp{Code: err.Error()}, nil)
	}
}

func makeStderr(dir string, id int32) string {
	return path.Join(dir, strconv.Itoa(int(id)), "stderr.log")
}

func (er *Executor) onProcExec(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessExecReq)
	logger.GetSugar().Infof("onProcExec %v", msg)

	if p, ok := processCache[msg.GetId()]; ok && p.State == StateRunning {
		_ = replier.Reply(&protocol.ProcessExecResp{Pid: int32(p.Pid())}, nil)
		return
	}

	// 创建文件目录
	fileDir := path.Join(msg.GetDir(), msg.GetKey())
	_ = os.MkdirAll(path.Dir(fileDir), os.ModePerm)

	// 配置文件
	if len(msg.GetConfig()) > 0 {
		for name, ctx := range msg.GetConfig() {
			filename := path.Join(fileDir, name)
			_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
			_ = ioutil.WriteFile(filename, []byte(ctx), os.ModePerm)
		}
	}

	// 错误信息文件
	filename := path.Join(fileDir, "stderr.log")
	errFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
		return
	}

	ecmd := exec.Command(msg.GetName(), msg.GetArgs()...)
	ecmd.Dir = msg.GetDir()
	ecmd.Stderr = errFile

	if p, err := ProcessWithCmd(ecmd, func(process *Process) {
		_ = errFile.Close()
		if process.State == StateStopped {
			delete(processCache, process.ID)
		}
		saveCache()
	}); err != nil {
		_ = errFile.Close()
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
	} else {
		p.ID = msg.GetId()
		p.Key = msg.GetKey()
		p.Stderr = filename
		p.Command = ecmd.String()
		processCache[p.ID] = p
		saveCache()
		_ = replier.Reply(&protocol.ProcessExecResp{Pid: int32(p.Pid())}, nil)
	}
}

func (er *Executor) onProcSignal(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessSignalReq)
	logger.GetSugar().Infof("onProcSignal %v", msg)

	if err := syscall.Kill(int(msg.GetPid()), syscall.Signal(msg.GetSignal())); err != nil {
		_ = replier.Reply(&protocol.ProcessSignalResp{Code: err.Error()}, nil)
	} else {
		_ = replier.Reply(&protocol.ProcessSignalResp{}, nil)
	}
}

func (er *Executor) onProcState(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessStateReq)
	//logger.GetSugar().Infof("onProcState %v", msg)

	states := map[int32]*protocol.ProcessState{}
	for _, id := range msg.GetIds() {
		state := &protocol.ProcessState{
			State: StateStopped,
		}
		if p, ok := processCache[id]; ok {
			state.Pid = int32(p.Pid())
			state.State = p.State
			if p.State == StateExited {
				if data, err := ioutil.ReadFile(p.Stderr); err == nil {
					state.ExitMsg = string(data)
				}
			}
		}
		states[id] = state
	}
	_ = replier.Reply(&protocol.ProcessStateResp{States: states}, nil)

}

func (er *Executor) onLogFile(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.LogFileReq)
	//logger.GetSugar().Infof("onLogFile %v", msg)

	if file, err := os.Open(msg.GetFilename()); err == nil {
		payload := int64(msg.GetPayload())
		if payload > protocol.BuffSize-(protocol.HeadSize+100) {
			payload = protocol.BuffSize - (protocol.HeadSize + 100)
		}
		info, _ := file.Stat()
		size := info.Size()

		off := int64(0)
		if payload >= size {
			off = 0
			payload = size
		} else {
			off = size - payload
		}

		data := make([]byte, payload)
		if _, err := file.ReadAt(data, off); err == nil {
			_ = replier.Reply(&protocol.LogFileResp{Context: data}, nil)
			return
		}
	}
	_ = replier.Reply(&protocol.LogFileResp{Context: nil}, nil)
}
