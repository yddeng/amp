package exec

import (
	"amp/common"
	"amp/protocol"
	"bytes"
	"fmt"
	"github.com/yddeng/dnet/drpc"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func (er *Executor) onCmdExec(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.CmdExecReq)
	log.Printf("onCmdExec %v", msg)

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
	log.Printf("onProcExec id:%d name:%s dir:%s args:%v", msg.GetId(), msg.GetName(), msg.GetDir(), msg.GetArgs())

	if p, ok := watchProcess[msg.GetId()]; ok && p.GetState() == common.StateRunning {
		_ = replier.Reply(&protocol.ProcessExecResp{Pid: int32(p.Pid)}, nil)
		return
	}

	// 创建文件目录
	fileDir := path.Join(msg.GetDir(), common.AmpDir, msg.GetName())
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
		return
	}

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
	ecmd := exec.Command(msg.GetArgs()[0], msg.GetArgs()[1:]...)
	ecmd.Dir = msg.GetDir()
	ecmd.Stderr = errFile

	if err := ecmd.Start(); err != nil {
		log.Println("command start error", err)
		_ = errFile.Close()
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
		return
	}

	p, err := NewProcess(int32(ecmd.Process.Pid))
	if err != nil {
		log.Println("new process error", err)
		_ = errFile.Close()
		_ = replier.Reply(&protocol.ProcessExecResp{Code: err.Error()}, nil)
		return
	}

	p.ID = msg.GetId()
	p.Name = msg.GetName()
	p.Stderr = filename
	watchProcess[p.ID] = p
	saveProcess()
	_ = replier.Reply(&protocol.ProcessExecResp{Pid: int32(p.Pid)}, nil)

	p.watch(func(process *Process) {
		er.Submit(func() {
			_ = errFile.Close()
			if process.GetState() == common.StateStopped {
				delete(watchProcess, process.ID)
			}
			saveProcess()
		})
	})
}

func (er *Executor) onProcSignal(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessSignalReq)
	log.Printf("onProcSignal %v", msg)

	if err := syscall.Kill(int(msg.GetPid()), syscall.Signal(msg.GetSignal())); err != nil {
		_ = replier.Reply(&protocol.ProcessSignalResp{Code: err.Error()}, nil)
	} else {
		_ = replier.Reply(&protocol.ProcessSignalResp{}, nil)
	}
}

func (er *Executor) onProcState(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.ProcessStateReq)
	//log.Printf("onProcState %v", msg)

	states := map[int32]*protocol.ProcessState{}
	for _, id := range msg.GetIds() {
		state := &protocol.ProcessState{
			State: common.StateStopped,
		}
		if p, ok := watchProcess[id]; ok {
			state.Pid = int32(p.Pid)
			state.State = p.State
			if p.State == common.StateExited {
				if data, err := ioutil.ReadFile(p.Stderr); err == nil {
					state.ExitMsg = string(data)
				}
			} else if p.State == common.StateRunning {
				state.Cpu = fmt.Sprintf("%.1f%", p.CPUPercent())
				state.Mem = fmt.Sprintf("%.1f%", p.MemoryPercent())
			}
		}
		states[id] = state
	}
	//log.Println("onProcState", 222, states)
	if err := replier.Reply(&protocol.ProcessStateResp{States: states}, nil); err != nil {
		log.Println(err)
	}

}

func (er *Executor) onLogFile(replier *drpc.Replier, req interface{}) {
	msg := req.(*protocol.LogFileReq)
	//log.Printf("onLogFile %v", msg)

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
