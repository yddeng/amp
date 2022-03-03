package exec

import (
	"amp/common"
	"amp/util"
	psProc "github.com/shirou/gopsutil/process"
	"io/ioutil"
	"log"
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
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
	State      string `json:"state"`
	Pid        int32  `json:"pid"`
	Stderr     string `json:"stderr"`

	process *psProc.Process
	mu      sync.Mutex
	die     chan struct{}
	dieOnce sync.Once
}

func (this *Process) watch(callback func(process *Process)) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case <-this.die:
				ticker.Stop()
				return
			case <-ticker.C:
			}

			isRunning, err := this.process.IsRunning()
			if err != nil || !isRunning {
				ticker.Stop()
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
				return
			}
		}
	}()
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

/*
func ProcessWithCmd(cmd *exec.Cmd, callback func(process *Process)) (*Process, error) {
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	process := &Process{
		State: common.StateRunning,
		Cmd:   cmd,
		done:  make(chan struct{}),
		Pid:   cmd.Process.Pid,
	}

	go func() {
		var err error
		defer er.Submit(func() {
			if err != nil {
				log.Printf("process %d done err %v\n", process.Pid, err)
				if state, ok := err.(*exec.ExitError); ok {
					log.Printf("process %d exiterror %s\n", process.Pid, state.String())
					// !success
					if state.ProcessState.Exited() {
						// exit
						process.State = common.StateExited
					} else {
						// signal 人为操作，视为正常停机
						process.State = common.StateStopped
					}
				} else {
					// 异常退出
					log.Printf("process %d exit %v\n", process.Pid, err)
					process.State = common.StateExited
				}
			} else {
				// err == nil && success
				process.State = common.StateStopped
			}
			close(process.done)
			callback(process)
		})
		err = cmd.Wait()
	}()
	return process, nil
}

// 根据pid 绑定进程，能监听停止状态
func ProcessWithPid(pid int, callback func(process *Process)) (*Process, error) {
	// todo 总是会成功，需要修改
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}
	log.Println("ProcessWithPid", pid, p, err)
	process := &Process{
		State:   common.StateRunning,
		Process: p,
		done:    make(chan struct{}),
		Pid:     pid,
	}

	go func() {
		var state *os.ProcessState
		var err error
		defer er.Submit(func() {
			if err != nil {
				// 异常退出
				log.Printf("process %d exit %s\n", process.Pid, err)
				process.State = common.StateExited
			} else {
				if !state.Success() {
					if state.Exited() {
						// exit
						process.State = common.StateExited
					} else {
						// signal 人为操作，视为正常停机
						process.State = common.StateStopped
					}
				} else {
					// success code=0
					process.State = common.StateStopped
				}

			}
			close(process.done)
			callback(process)
		})
		state, err = p.Wait()
	}()

	return process, nil
}
*/

var (
	watchProcess = map[int32]*Process{}
	processFile  string
)

func loadProcess(dataPath string) {
	var processMap map[int32]*Process
	processFile = path.Join(dataPath, "exec_info.json")
	if err := util.DecodeJsonFromFile(&processMap, processFile); err == nil {
		for _, p := range processMap {
			if p.GetState() != common.StateRunning {
				watchProcess[p.ID] = p
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
					watchProcess[proc.ID] = proc
					proc.watch(func(process *Process) {
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
	if err := util.EncodeJsonToFile(watchProcess, processFile); err != nil {
		log.Printf("saveProcess faield %v", err)
	}
}
