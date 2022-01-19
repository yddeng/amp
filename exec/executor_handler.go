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

func (er *Executor) onProcExec(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessExecReq)
	logger.GetSugar().Infof("onProcExec %v", msg)

	if len(msg.GetConfig()) > 0 {
		for name, ctx := range msg.GetConfig() {
			_ = os.MkdirAll(path.Dir(name), os.ModePerm)
			_ = ioutil.WriteFile(name, []byte(ctx), os.ModePerm)
		}
	}

	ecmd := exec.Command(msg.GetName(), msg.GetArgs()...)
	ecmd.Dir = msg.GetDir()

	// 错误信息重定向
	var errFile *os.File
	var err error
	if msg.GetStderr() != "" {
		_ = os.MkdirAll(path.Dir(msg.GetStderr()), os.ModePerm)
		errFile, err = os.OpenFile(msg.GetStderr(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
			return
		}
		ecmd.Stderr = errFile
	}

	cmd := CommandWithCmd(ecmd)
	if err = cmd.Run(0, func(cmd *Cmd, err error) {
		er.Submit(func() {
			if errFile != nil {
				_ = errFile.Close()
			}
		})
	}); err != nil {
		if errFile != nil {
			_ = errFile.Close()
		}
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
	} else {
		_ = replier.Reply(&protocol.ProcessExecResp{Pid: int32(cmd.Pid())}, nil)
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

func (er *Executor) onProcIsAlive(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessIsAliveReq)
	logger.GetSugar().Infof("onProcIsAlive %v", msg)

	ret := map[int32]bool{}
	for id, pid := range msg.GetPid() {
		if err := syscall.Kill(int(pid), 0); err == nil {
			ret[id] = true
		} else {
			ret[id] = false
		}
	}

	_ = replier.Reply(&protocol.ProcessIsAliveResp{Alive: ret}, nil)

}

func (er *Executor) onLogFile(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.LogFileReq)
	logger.GetSugar().Infof("onLogFile %v", msg)

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
