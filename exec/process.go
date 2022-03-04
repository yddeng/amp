package exec

import (
	"amp/common"
	"amp/util"
	psProc "github.com/shirou/gopsutil/process"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
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
	Pid        int32  `json:"pid"`
	Stderr     string `json:"stderr"`
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	CreateTime int64  `json:"create_time"`

	process *psProc.Process
	mu      sync.Mutex
}

func (this *Process) waitCmd(cmd *exec.Cmd, callback func(process *Process)) {
	go func() {
		err := cmd.Wait()

		this.mu.Lock()
		if err != nil {
			if state, ok := err.(*exec.ExitError); ok {
				// !success
				if state.ProcessState.Exited() {
					// exit
					this.State = common.StateExited
				} else {
					// signal 人为操作，视为正常停机
					this.State = common.StateStopped
				}
			} else {
				// 异常退出
				this.State = common.StateExited
			}
		} else {
			// err == nil && success
			this.State = common.StateStopped
		}
		this.mu.Unlock()
		callback(this)
	}()
}

func (this *Process) waitChild(proc *os.Process, callback func(process *Process)) {
	go func() {
		state, err := proc.Wait()
		this.mu.Lock()
		if err != nil {
			// 异常退出
			this.State = common.StateExited
		} else {
			if !state.Success() {
				if state.Exited() {
					// exit
					this.State = common.StateExited
				} else {
					// signal 人为操作，视为正常停机
					this.State = common.StateStopped
				}
			} else {
				// success code=0
				this.State = common.StateStopped
			}

		}
		this.mu.Unlock()
		callback(this)

	}()
}

func (this *Process) waitNoChild(callback func(process *Process)) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			<-ticker.C

			isRunning, err := this.process.IsRunning()
			if err != nil || !isRunning {
				if this.Stderr != "" {
					data, err := ioutil.ReadFile(this.Stderr)
					this.mu.Lock()
					if err == nil && len(data) != 0 {
						this.State = common.StateExited
					} else {
						this.State = common.StateStopped
					}
					this.mu.Unlock()
				} else {
					this.mu.Lock()
					this.State = common.StateStopped
					this.mu.Unlock()
				}
				callback(this)
				ticker.Stop()
				return
			}
		}
	}()
}

func (this *Process) wait(callback func(process *Process)) error {
	pp, err := this.process.Parent()
	if err != nil {
		return err
	} else {
		if pp.Pid == int32(os.Getpid()) {
			ppp, err := os.FindProcess(int(this.Pid))
			if err != nil {
				return err
			}
			this.waitChild(ppp, callback)
		} else {
			this.waitNoChild(callback)
		}
	}
	return nil
}

func (this *Process) GetState() string {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.State
}

func (this *Process) CPUPercent() float64 {
	if this.process == nil {
		return 0
	}
	percent, err := this.process.CPUPercent()
	if err != nil {
		return 0
	}
	return percent
}

func (this *Process) MemoryPercent() float32 {
	if this.process == nil {
		return 0
	}
	percent, err := this.process.MemoryPercent()
	if err != nil {
		return 0
	}
	return percent
}

func NewProcess(pid int32) (*Process, error) {
	p, err := psProc.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	createTime, err := p.CreateTime()
	if err != nil {
		return nil, err
	}

	this := &Process{
		CreateTime: createTime,
		State:      common.StateRunning,
		Pid:        pid,
		process:    p,
	}
	return this, nil
}

var (
	waitProcess = map[int32]*Process{}
	processFile string
)

func loadProcess(dataPath string) {
	var processMap map[int32]*Process
	processFile = path.Join(dataPath, "exec_info.json")
	if err := util.DecodeJsonFromFile(&processMap, processFile); err == nil {
		for _, p := range processMap {
			if p.GetState() != common.StateRunning {
				waitProcess[p.ID] = p
				continue
			}

			if proc, err := NewProcess(p.Pid); err != nil {
				log.Printf("loadProcess %s faield %d %v", p.Name, p.Pid, err)
			} else {
				if p.CreateTime != proc.CreateTime {
					log.Printf("loadProcess %s faield %d create time not equal", p.Name, p.Pid)
				} else {
					proc.ID = p.ID
					proc.Name = p.Name
					proc.Stderr = p.Stderr
					waitProcess[proc.ID] = proc
					proc.waitNoChild(func(process *Process) {
						er.Submit(func() {
							saveProcess()
						})
					})
				}
			}
		}
	}
	saveProcess()
}

func saveProcess() {
	if err := util.EncodeJsonToFile(waitProcess, processFile); err != nil {
		log.Printf("saveProcess faield %v", err)
	}
}
