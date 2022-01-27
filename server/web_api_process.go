package server

import (
	"amp/exec"
	"amp/protocol"
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
	isAutoStart    bool
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
	for nav := range processMgr.Groups {
		if strings.HasPrefix(nav, req.Group) {
			delg[nav] = struct{}{}
		}
	}

	if len(delg) == 0 {
		done.result.Message = "不存在的分组"
		return
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
				if strings.HasPrefix(nav, req.Group) {
					s[v.ID] = v
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
		Status: exec.StateStopped,
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

	for _, g := range req.Groups {
		if _, ok := processMgr.Groups[g]; !ok {
			this.createGroupPath(g)
		}
	}

	p, ok := processMgr.Process[req.ID]
	if !ok || !(p.State.Status == exec.StateStopped ||
		p.State.Status == exec.StateExited) {
		done.result.Message = "当前状态不允许修改"
		return
	}

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
	if !ok || !(p.State.Status == exec.StateStopped ||
		p.State.Status == exec.StateExited) {
		done.result.Message = "当前状态不允许删除"
		return
	}

	delete(processMgr.Process, req.ID)
	saveStore(snProcessMgr)
}

func processTick() {
	rpcReq := map[string]*protocol.ProcessStateReq{}
	for _, p := range processMgr.Process {
		if p.State.Status == exec.StateStarting ||
			p.State.Status == exec.StateRunning ||
			p.State.Status == exec.StateStopping {
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
			continue
		}
		_ = center.Go(node, req, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
			if e != nil {
				return
			}
			change := false
			rpcResp := i.(*protocol.ProcessStateResp)
			for id, state := range rpcResp.GetStates() {
				p, ok := processMgr.Process[int(id)]
				if !ok {
					continue
				}
				pState := p.State.Status
				if pState != state.GetState() {
					if pState == exec.StateStarting {
						p.State.Status = state.GetState()
						if state.GetState() == exec.StateRunning {
							p.State.Pid = state.Pid
						} else if state.GetState() == exec.StateExited {
							p.State.ExitMsg = state.GetExitMsg()
							p.State.isAutoStart = true
						} else {
							// Stopped
						}
						change = true
					} else if pState == exec.StateRunning {
						p.State.Status = state.GetState()
						if state.GetState() == exec.StateExited {
							p.State.ExitMsg = state.GetExitMsg()
							p.State.isAutoStart = true
						} else {
							// Stopped
						}
						change = true
					} else {
						// Stopping
						if state.GetState() == exec.StateRunning {
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
							p.State.Status = exec.StateStopped
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
		if p.State.isAutoStart &&
			p.State.Status == exec.StateExited &&
			p.State.AutoStartTimes < p.AutoStartTimes {

			node, ok := nodes[p.Node]
			if !ok || !node.Online() {
				continue
			}

			log.Printf("process %d auto start times %d\n", p.ID, p.State.AutoStartTimes)
			if err := p.start(node, func(code string, err error) {}); err == nil {
				p.State.AutoStartTimes += 1
				p.State.Status = exec.StateStarting
				p.State.Timestamp = NowUnix()
				p.State.ExitMsg = ""
				saveStore(snProcessMgr)
			}
		}
	}
}

//func makePath(dir string, id int) string {
//	return path.Join(dir, strconv.Itoa(id))
//}

func (p *Process) start(node *Node, callback func(code string, err error)) error {
	key := p.Name
	configs := make(map[string]string, len(p.Config))
	for _, cfg := range p.Config {
		configs[cfg.Name] = cfg.Context
	}

	cmd := strings.ReplaceAll(p.Command, "{{path}}", key)
	cmds := strings.Split(cmd, " ")
	rpcReq := &protocol.ProcessExecReq{
		Id:     int32(p.ID),
		Key:    key,
		Dir:    p.Dir,
		Name:   cmds[0],
		Args:   cmds[1:],
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
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok ||
		p.State.Status == exec.StateStarting ||
		p.State.Status == exec.StateRunning ||
		p.State.Status == exec.StateStopping {
		done.result.Message = "当前状态不允许启动"
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
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
			Status:    exec.StateStarting,
			Timestamp: NowUnix(),
		}
		saveStore(snProcessMgr)
	}
}

func (*processHandler) Stop(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok || p.State.Status != exec.StateRunning {
		done.result.Message = "当前状态不允许停止"
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
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
		p.State.Status = exec.StateStopping
		p.State.Timestamp = NowUnix()
		saveStore(snProcessMgr)
	}
}
