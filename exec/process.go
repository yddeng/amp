package exec

import (
	"amp/common"
	"amp/util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
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
	// find by pid
	Process *os.Process `json:"_"`
	// Cmd
	Cmd *exec.Cmd `json:"_"`

	done chan struct{}

	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Command string `json:"command"`
	State   string `json:"state"`
	Pid     int    `json:"pid"`
	Stderr  string `json:"stderr"`
}

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

func (this *Process) Signal(sig syscall.Signal) error {
	return syscall.Kill(this.Pid, sig)
}

func (this *Process) Done() bool {
	select {
	case <-this.done:
		return true
	default:
		return false
	}
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
	var processMap map[int32]*Process
	cacheFile = path.Join(dataPath, "exec_info.json")
	if err := util.DecodeJsonFromFile(&processMap, cacheFile); err == nil {
		for _, p := range processMap {
			if m, err := ProcessCollect(p.Pid); err != nil {
				log.Printf("loadCache %s faield %d %v", p.Name, p.Pid, err)
			} else {
				if p.Command != m.Args {
					log.Printf("loadCache %s faield %d command not equal  %s <=> %s", p.Name, p.Pid, p.Command, m.Args)
				} else {
					log.Printf("loadCache %s process %d ok", p.Name, p.Pid)
					p.done = make(chan struct{})
					processCache[p.ID] = p
					go func(p *Process) {
						ticker := time.NewTicker(time.Second)
						for {
							select {
							case <-p.done:
								ticker.Stop()
								return
							case <-ticker.C:
								if !p.IsAlive() {
									data, err := ioutil.ReadFile(p.Stderr)
									if err == nil && len(data) != 0 {
										p.State = common.StateExited
									} else {
										p.State = common.StateStopped
									}
									close(p.done)
								}
							}
						}
					}(p)
				}
			}
		}
	}
	saveCache()
}

func saveCache() {
	if err := util.EncodeJsonToFile(processCache, cacheFile); err != nil {
		log.Printf("saveCache faield %v", err)
	}
}
