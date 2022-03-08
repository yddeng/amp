package server

import (
	"amp/back-go/common"
	"amp/back-go/protocol"
	"fmt"
	"github.com/yddeng/dnet/drpc"
	"log"
	"strings"
	"syscall"
)

type ProcessConfig struct {
	Name    string `json:"name"`
	Context string `json:"context"`
}

type ProcessState struct {
	// 状态
	Status string `json:"status"`
	// 运行时Pid， Running、Stopping 时可用，启动时重置
	Pid int32 `json:"pid"`
	// 时间戳 秒， 启动、停止时设置
	Timestamp int64 `json:"timestamp"`
	// Exited 信息，启动时重置
	ExitMsg string `json:"exit_msg"`
	// 已经自动重启次数，启动时重置
	AutoStartTimes int `json:"auto_start_times"`

	Cpu float64 `json:"cpu"`
	Mem float64 `json:"mem"`
}

type Process struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Dir      string           `json:"dir"`
	Config   []*ProcessConfig `json:"config"`
	Command  string           `json:"command"`
	Groups   []string         `json:"groups"`
	Node     string           `json:"node"`
	User     string           `json:"user"`
	CreateAt int64            `json:"create_at"`

	State ProcessState `json:"state"`

	// 子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭
	// 默认值为999 。。非必须设置
	Priority int `json:"priority"`
	// 子进程启动多少秒之后，此时状态如果是running，则我们认为启动成功了
	// 默认值为2 。。非必须设置
	StartSecs int64 `json:"start_secs"`
	// 这个是当我们向子进程发送stop信号后，到系统返回信息所等待的最大时间。
	// 超过这个时间会向该子进程发送一个强制kill的信号。
	// 默认为10秒。。非必须设置
	StopWaitSecs int64 `json:"stop_wait_secs"`
	// 进程状态为 Exited时，自动重启
	// 默认为3次。。非必须设置
	AutoStartTimes int `json:"auto_start_times"`
}

func processTick() {
	rpcReq := map[string]*protocol.ProcessStateReq{}
	for _, p := range processMgr.Process {
		if !(p.State.Status == common.StateStopped ||
			p.State.Status == common.StateExited) {
			req, ok := rpcReq[p.Node]
			if !ok {
				req = &protocol.ProcessStateReq{
					Ids: make([]int32, 0, 4),
				}
				rpcReq[p.Node] = req
			}
			req.Ids = append(req.Ids, int32(p.ID))
		}
	}

	for n, req := range rpcReq {
		node, ok := nodes[n]
		if !ok || !node.Online() {
			// 节点不在线 设置状态为 unknown
			change := false
			for _, id := range req.GetIds() {
				p, ok := processMgr.Process[int(id)]
				if !ok {
					continue
				}
				if p.State.Status != common.StateUnknown {
					p.State.Status = common.StateUnknown
					change = true
				}
			}
			if change {
				saveStore(snProcessMgr)
			}
			continue
		}
		_ = center.Go(node, req, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
			if e != nil {
				return
			}
			change := false
			rpcResp := i.(*protocol.ProcessStateResp)
			//log.Println(22, rpcResp)
			for id, state := range rpcResp.GetStates() {
				p, ok := processMgr.Process[int(id)]
				if !ok {
					continue
				}

				p.State.Pid = state.Pid
				p.State.Mem = state.GetMem()
				p.State.Cpu = state.GetCpu()
				switch p.State.Status {
				case common.StateUnknown:
					p.State.Status = state.GetState()
					switch state.GetState() {
					case common.StateRunning:
					case common.StateStopped:
					case common.StateExited:
						p.State.AutoStartTimes = p.AutoStartTimes // 未知状态不重启
						p.State.ExitMsg = state.GetExitMsg()
					}
					change = true
				case common.StateStarting:
					switch state.GetState() {
					case common.StateRunning:
						if NowUnix() >= p.State.Timestamp+p.StartSecs {
							// 启动时间
							p.State.Status = common.StateRunning
							change = true
						}
					case common.StateStopped:
						p.State.Status = common.StateStopped
						change = true
					case common.StateExited:
						p.State.Status = common.StateExited
						p.State.AutoStartTimes = p.AutoStartTimes // 启动阶段不重启
						p.State.ExitMsg = state.GetExitMsg()
						change = true
					}
				case common.StateRunning:
					switch state.GetState() {
					case common.StateRunning:
					case common.StateStopped:
						p.State.Status = common.StateStopped
						change = true
					case common.StateExited:
						p.State.Status = common.StateExited
						p.State.ExitMsg = state.GetExitMsg()
						change = true
					}
				case common.StateStopping:
					switch state.GetState() {
					case common.StateRunning:
						if NowUnix() >= p.State.Timestamp+p.StopWaitSecs {
							// 停止时间超时 ，强行停止
							_ = center.Go(node, &protocol.ProcessSignalReq{
								Pid:    p.State.Pid,
								Signal: int32(syscall.SIGKILL),
							}, drpc.DefaultRPCTimeout, func(i interface{}, e error) {})
						}
					case common.StateStopped:
						p.State.Status = common.StateStopped
						change = true
					case common.StateExited:
						p.State.Status = common.StateExited
						p.State.AutoStartTimes = p.AutoStartTimes // 停止阶段不重启
						p.State.ExitMsg = state.GetExitMsg()
						change = true
					}
				}
			}
			if change {
				saveStore(snProcessMgr)
			}
		})
	}
}

func processAutoStart() {
	for _, p := range processMgr.Process {
		if p.State.Status == common.StateExited &&
			p.State.AutoStartTimes < p.AutoStartTimes {

			node, ok := nodes[p.Node]
			if !ok || !node.Online() {
				continue
			}

			log.Printf("process %d auto start times %d\n", p.ID, p.State.AutoStartTimes)
			if err := p.start(node, func(code string, err error) {}); err == nil {
				p.State.AutoStartTimes += 1
				p.State.Status = common.StateStarting
				p.State.Timestamp = NowUnix()
				p.State.ExitMsg = ""
				saveStore(snProcessMgr)
			}
		}
	}
}

func (p *Process) start(node *Node, callback func(code string, err error)) error {
	configs := make(map[string]string, len(p.Config))
	for _, cfg := range p.Config {
		configs[cfg.Name] = cfg.Context
	}

	cmd := strings.ReplaceAll(p.Command, "{{path}}", fmt.Sprintf("%s/%s", common.AmpDir, p.Name))
	rpcReq := &protocol.ProcessExecReq{
		Id:     int32(p.ID),
		Dir:    p.Dir,
		Name:   p.Name,
		Args:   strings.Fields(cmd),
		Config: configs,
	}

	return center.Go(node, rpcReq, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
		if e != nil {
			callback("", e)
			return
		}
		rpcResp := i.(*protocol.ProcessExecResp)
		callback(rpcResp.GetCode(), nil)
	})
}
