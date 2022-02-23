package server

import (
	"amp/common"
	"amp/protocol"
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

	Cpu float64 `json:"_"`
	Mem float64 `json:"_"`
}

type Process struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Dir      string           `json:"dir"`
	Config   []*ProcessConfig `json:"config"`
	Command  string           `json:"command"`
	State    ProcessState     `json:"state"`
	Groups   []string         `json:"groups"`
	Node     string           `json:"node"`
	User     string           `json:"user"`
	CreateAt int64            `json:"create_at"`

	// 子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭
	// 默认值为999 。。非必须设置
	Priority int `json:"priority"`
	// 子进程启动多少秒之后，此时状态如果是running，则我们认为启动成功了
	// 默认值为1 。。非必须设置
	//StartSecs int `json:"start_secs"`
	// 这个是当我们向子进程发送stop信号后，到系统返回信息所等待的最大时间。
	// 超过这个时间会向该子进程发送一个强制kill的信号。
	// 默认为10秒。。非必须设置
	StopWaitSecs int `json:"stop_wait_secs"`
	// 进程状态为 Exited时，自动重启
	// 默认为3次。。非必须设置
	AutoStartTimes int `json:"auto_start_times"`
}

type ProcessMgr struct {
	GenID   int                 `json:"gen_id"`
	Process map[int]*Process    `json:"process"`
	Groups  map[string]struct{} `json:"groups"` // 程序组 'nav1/nav2'
}

type processHandler struct {
}

func (*processHandler) GroupList(done *Done, user string) {
	//log.Printf("%s by(%s) \n", done.route, user)
	defer func() { done.Done() }()
	done.result.Data = processMgr.Groups
}

func (*processHandler) createGroupPath(group string) {
	// 创建路径
	ss := strings.Split(group, "/")
	var nav string
	for i := 1; i <= len(ss); i++ {
		nav = strings.Join(ss[0:i], "/")
		if _, ok := processMgr.Groups[nav]; !ok {
			processMgr.Groups[nav] = struct{}{}
		}
	}
	saveStore(snProcessMgr)
}

func (this *processHandler) GroupAdd(done *Done, user string, req struct {
	Group string `json:"group"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := processMgr.Groups[req.Group]; !ok {
		this.createGroupPath(req.Group)
	}
}

func (*processHandler) GroupRemove(done *Done, user string, req struct {
	Group string `json:"group"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	delg := map[string]struct{}{}
	if _, ok := processMgr.Groups[req.Group]; ok {
		delg[req.Group] = struct{}{}
	} else {
		done.result.Message = "不存在的分组"
		return
	}

	// 移除子分组
	// 防止前缀一致，加上分割符， all、 alltest 移除 all
	prefix := req.Group + "/"
	for nav := range processMgr.Groups {
		if strings.HasPrefix(nav, prefix) {
			delg[nav] = struct{}{}
		}
	}

	for _, v := range processMgr.Process {
		for _, g := range v.Groups {
			if _, ok := delg[g]; ok {
				done.result.Message = "当前分组还存在进程，不允许删除"
				return
			}
		}
	}

	for name := range delg {
		delete(processMgr.Groups, name)
	}
	saveStore(snProcessMgr)
}

func (*processHandler) List(done *Done, user string, req struct {
	Group string `json:"group"`
}) {
	//log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if req.Group == "" {
		done.result.Data = processMgr.Process
	} else {
		s := make(map[int]*Process, len(processMgr.Process))
		for _, v := range processMgr.Process {
			for _, nav := range v.Groups {
				if nav == req.Group {
					s[v.ID] = v
				} else {
					prefix := req.Group + "/"
					if strings.HasPrefix(nav, prefix) {
						s[v.ID] = v
					}
				}
			}
		}
		done.result.Data = s
	}
}

func (this *processHandler) Create(done *Done, user string, req struct {
	Name           string           `json:"name"`
	Dir            string           `json:"dir"`
	Config         []*ProcessConfig `json:"config"`
	Command        string           `json:"command"`
	Groups         []string         `json:"groups"`
	Node           string           `json:"node"`
	Priority       int              `json:"priority"`
	StopWaitSecs   int              `json:"stop_wait_secs"`
	AutoStartTimes int              `json:"auto_start_times"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	for _, p := range processMgr.Process {
		if p.Name == req.Name {
			done.result.Message = "程序名重复"
			return
		}
	}

	for _, g := range req.Groups {
		if _, ok := processMgr.Groups[g]; !ok {
			this.createGroupPath(g)
		}
	}

	processMgr.GenID++
	id := processMgr.GenID
	p := new(Process)
	p.ID = id
	p.Name = req.Name
	p.Dir = req.Dir
	p.Config = req.Config
	p.Command = req.Command
	p.State = ProcessState{
		Status: common.StateStopped,
	}
	p.Groups = req.Groups
	p.Node = req.Node
	p.User = user
	p.CreateAt = NowUnix()
	p.Priority = req.Priority
	p.StopWaitSecs = req.StopWaitSecs
	p.AutoStartTimes = req.AutoStartTimes

	processMgr.Process[id] = p
	saveStore(snProcessMgr)
}

func (this *processHandler) Update(done *Done, user string, req struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Dir            string           `json:"dir"`
	Config         []*ProcessConfig `json:"config"`
	Command        string           `json:"command"`
	Groups         []string         `json:"groups"`
	Node           string           `json:"node"`
	Priority       int              `json:"priority"`
	StopWaitSecs   int              `json:"stop_wait_secs"`
	AutoStartTimes int              `json:"auto_start_times"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	for _, p := range processMgr.Process {
		if p.Name == req.Name && p.ID != req.ID {
			done.result.Message = "程序名重复"
			return
		}
	}

	for _, g := range req.Groups {
		if _, ok := processMgr.Groups[g]; !ok {
			this.createGroupPath(g)
		}
	}

	p, ok := processMgr.Process[req.ID]
	if !ok || !(p.State.Status == common.StateStopped ||
		p.State.Status == common.StateExited) {
		done.result.Message = "当前状态不允许修改"
		return
	}

	p.Name = req.Name
	p.Dir = req.Dir
	p.Config = req.Config
	p.Command = req.Command
	p.Groups = req.Groups
	p.Node = req.Node
	p.Priority = req.Priority
	p.StopWaitSecs = req.StopWaitSecs
	p.AutoStartTimes = req.AutoStartTimes

	saveStore(snProcessMgr)
}

func (*processHandler) Delete(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok || !(p.State.Status == common.StateStopped ||
		p.State.Status == common.StateExited) {
		done.result.Message = "当前状态不允许删除"
		return
	}

	delete(processMgr.Process, req.ID)
	saveStore(snProcessMgr)
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
				//log.Println("11", e)
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
				p.State.Mem = state.GetMem()
				p.State.Cpu = state.GetCpu()
				pState := p.State.Status
				if pState != state.GetState() {
					if pState == common.StateUnknown || pState == common.StateStarting {
						p.State.Status = state.GetState()
						if state.GetState() == common.StateRunning {
							p.State.Pid = state.Pid
						} else if state.GetState() == common.StateExited {
							// 启动报错不重启
							p.State.AutoStartTimes = p.AutoStartTimes
							p.State.ExitMsg = state.GetExitMsg()
						} else {
							// Stopped
						}
						change = true
					} else if pState == common.StateRunning {
						p.State.Status = state.GetState()
						if state.GetState() == common.StateExited {
							p.State.ExitMsg = state.GetExitMsg()
						} else {
							// Stopped
						}
						change = true
					} else {
						// Stopping
						if state.GetState() == common.StateRunning {
							subUnix := NowUnix() - p.State.Timestamp
							if int(subUnix) >= p.StopWaitSecs {
								_ = center.Go(node, &protocol.ProcessSignalReq{
									Pid:    p.State.Pid,
									Signal: int32(syscall.SIGKILL),
								}, drpc.DefaultRPCTimeout, func(i interface{}, e error) {})
							}
						} else {
							// Stopped || Exited
							p.State.ExitMsg = state.GetExitMsg()
							p.State.Status = common.StateStopped
							change = true
						}
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

func (*processHandler) Start(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	p, ok := processMgr.Process[req.ID]
	if !ok ||
		!(p.State.Status == common.StateStopped ||
			p.State.Status == common.StateExited) {
		done.result.Message = "当前状态不允许启动"
		done.Done()
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
		done.Done()
		return
	}

	if err := p.start(node, func(code string, err error) {
		if err == nil {
			done.result.Message = code
		}
		done.Done()
	}); err != nil {
		done.result.Message = err.Error()
		done.Done()
	} else {
		p.State = ProcessState{
			Status:    common.StateStarting,
			Timestamp: NowUnix(),
		}
		saveStore(snProcessMgr)
	}
}

func (*processHandler) Stop(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	p, ok := processMgr.Process[req.ID]
	if !ok || p.State.Status != common.StateRunning {
		done.result.Message = "当前状态不允许停止"
		done.Done()
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
		done.Done()
		return
	}

	rpcReq := &protocol.ProcessSignalReq{
		Pid:    p.State.Pid,
		Signal: int32(syscall.SIGTERM),
	}

	if err := center.Go(node, rpcReq, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
		if e != nil {
			done.Done()
			return
		}
		rpcResp := i.(*protocol.ProcessSignalResp)
		if rpcResp.GetCode() != "" {
			done.result.Message = rpcResp.GetCode()
		}
		done.Done()
	}); err != nil {
		done.result.Message = err.Error()
		done.Done()
	} else {
		p.State.Status = common.StateStopping
		p.State.Timestamp = NowUnix()
		saveStore(snProcessMgr)
	}
}
