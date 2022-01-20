package exec

import (
	"initial-server/util"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const (
	StateStarting = "Starting"
	StateRunning  = "Running"
	StateStopping = "Stopping"
	StateStopped  = "Stopped"
	StateExited   = "Exited"
)

type Daemon struct {
	// 这个选项是进程启动多少秒之后，此时状态如果是running，则我们认为启动成功了
	// 默认值为1 。。非必须设置
	StartSecs int
	// 当进程启动失败后，最大尝试启动的次数。。当超过预定次数后，将把此进程的状态置为Exited
	// 默认值为3 。。非必须设置。
	StartRetries int
	// 这个是当我们向子进程发送stop信号后，到系统返回信息所等待的最大时间。
	// 超过这个时间会向该子进程发送一个强制kill的信号。
	StopWaitSecs int
	// 运行状态
	State string
	// Cmd
	Cmd *exec.Cmd
}

type Process struct {
	// find by pid
	Process *os.Process `json:"_"`
	// state
	ProcessState *os.ProcessState `json:"_"`
	// Cmd
	Cmd *exec.Cmd `json:"_"`

	ID      int32  `json:"id"`
	Command string `json:"command"`
	State   string `json:"state"`
	Stderr  string `json:"stderr"`
}

func ProcessWithCmd(cmd *exec.Cmd, callback func(process *Process)) (*Process, error) {
	process := &Process{Cmd: cmd}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	process.State = StateRunning

	go func() {
		var err error
		defer er.Submit(func() {
			process.ProcessState = cmd.ProcessState
			if err != nil {
				if process.ProcessState.Exited() {
					// exit
					process.State = StateExited
				} else {
					// signal 人为操作，视为正常停机
					process.State = StateStopped
				}
			} else {
				// success code=0
				process.State = StateStopped
			}
			callback(process)
		})
		err = cmd.Wait()
	}()
	return process, nil
}

// 根据pid 绑定进程，能监听停止状态
func ProcessWithPid(pid int, callback func(process *Process)) (*Process, error) {
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}
	process := &Process{State: StateRunning, Process: p}

	go func() {
		var state *os.ProcessState
		//var err error
		defer er.Submit(func() {
			process.ProcessState = state
			if !state.Success() {
				if state.Exited() {
					// exit
					process.State = StateExited
				} else {
					// signal 人为操作，视为正常停机
					process.State = StateStopped
				}
			} else {
				// success code=0
				process.State = StateStopped
			}
			callback(process)
		})
		state, _ = p.Wait()
	}()

	return process, nil
}

func (this *Process) Pid() int {
	if this.Process == nil {
		return this.Cmd.Process.Pid
	} else {
		return this.Process.Pid
	}
}

func (this *Process) Signal(sig syscall.Signal) error {
	if this.Process == nil {
		return this.Cmd.Process.Signal(sig)
	} else {
		return this.Process.Signal(sig)
	}
}

func (this *Process) Done() bool {
	if this.ProcessState != nil {
		return true
	}
	return false
}

func (this *Process) IsAlive() bool {
	if this.Done() {
		return false
	} else {
		return this.Signal(syscall.Signal(0)) == nil
	}
}

var (
	processCache = map[int32]*Process{}
	cacheFile    string
)

func loadCache(dataPath string) {
	cacheFile = path.Join(dataPath, "exec_info.json")
	if err := util.DecodeJsonFromFile(&processCache, cacheFile); err == nil {
	}
}

func saveCache() {
	if err := util.EncodeJsonToFile(processCache, cacheFile); err != nil {

	}
}
