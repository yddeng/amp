package server

import (
	"github.com/yddeng/dnet/drpc"
	"initial-server/protocol"
	"log"
	"path"
	"strconv"
	"strings"
	"syscall"
)

type ProcessConfig struct {
	Name    string `json:"name"`
	Context string `json:"context"`
}

const (
	processStarting = "Starting"
	processRunning  = "Running"
	processFatal    = "Fatal"  // 设置重启时，超过重启次数
	processExited   = "Exited" // 程序退出
	processStopping = "Stopping"
	processStopped  = "Stopped"
)

type ProcessState struct {
	Pid       int32  `json:"pid"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

type Process struct {
	ID           int                 `json:"id"`
	Name         string              `json:"name"`
	Dir          string              `json:"dir"`
	Config       []*ProcessConfig    `json:"config"`
	Command      string              `json:"command"`
	Priority     int                 `json:"priority"`       // 子进程启动关闭优先级，优先级低的，最先启动，关闭的时候最后关闭	默认值为999 。。非必须设置
	StartRetries int                 `json:"start_retries"`  // 当进程启动失败后，最大尝试启动的次数。。当超过3次后，supervisor将把 此进程的状态置为FAIL	默认值为3 。。非必须设置
	StopWaitSecs int                 `json:"stop_wait_secs"` // 这个是当我们向子进程发送stopsignal信号后，到系统返回信息	给supervisord，所等待的最大时间。 超过这个时间，supervisord会向该	子进程发送一个强制kill的信号。
	State        ProcessState        `json:"state"`
	Groups       map[string]struct{} `json:"groups"`
	Node         string              `json:"node"`
	User         string              `json:"user"`
	CreateAt     int64               `json:"create_at"`
}

type ProcessMgr struct {
	GenID   int                 `json:"gen_id"`
	Process map[int]*Process    `json:"process"`
	Groups  map[string]struct{} `json:"groups"` // 程序组 'nav1/nav2'
}

type processHandler struct {
}

func (*processHandler) GroupList(done *Done, user string) {
	log.Printf("%s by(%s) \n", done.route, user)
	defer func() { done.Done() }()
	done.result.Data = processMgr.Groups
}

func (*processHandler) GroupAdd(done *Done, user string, req struct {
	Group string `json:"group"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := processMgr.Groups[req.Group]; !ok {
		// 创建路径
		ss := strings.Split(req.Group, "/")
		var nav string
		for i := 1; i <= len(ss); i++ {
			nav = strings.Join(ss[0:i], "/")
			if _, ok := processMgr.Groups[nav]; !ok {
				processMgr.Groups[nav] = struct{}{}
			}
		}
		saveStore(snProcessMgr)
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
		for g := range v.Groups {
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
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make(map[int]*Process, len(processMgr.Process))
	for _, v := range processMgr.Process {
		for nav := range v.Groups {
			if strings.HasPrefix(nav, req.Group) {
				s[v.ID] = v
			}
		}
	}
	done.result.Data = s
}

func (*processHandler) Create(done *Done, user string, req struct {
	Name         string           `json:"name"`
	Dir          string           `json:"dir"`
	Config       []*ProcessConfig `json:"config"`
	Command      string           `json:"command"`
	Priority     int              `json:"priority"`
	StartRetries int              `json:"start_retries"`
	StopWaitSecs int              `json:"stop_wait_secs"`
	Groups       []string         `json:"groups"`
	Node         string           `json:"node"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	gs := map[string]struct{}{}
	for _, g := range req.Groups {
		if _, ok := processMgr.Groups[g]; !ok {
			// 创建路径
			ss := strings.Split(g, "/")
			var nav string
			for i := 0; i <= len(ss); i++ {
				nav = strings.Join(ss[0:i], "/")
				if _, ok := processMgr.Groups[nav]; !ok {
					processMgr.Groups[nav] = struct{}{}
				}
			}
		}
		gs[g] = struct{}{}
	}

	processMgr.GenID++
	id := processMgr.GenID
	p := new(Process)
	p.ID = id
	p.Name = req.Name
	p.Dir = req.Dir
	p.Config = req.Config
	p.Command = req.Command
	p.Priority = req.Priority
	p.StartRetries = req.StartRetries
	p.StopWaitSecs = req.StartRetries
	p.State = ProcessState{
		Status: processStopped,
	}
	p.Groups = gs
	p.Node = req.Node
	p.User = user
	p.CreateAt = NowUnix()

	processMgr.Process[id] = p
	saveStore(snProcessMgr)
}

func (*processHandler) Update(done *Done, user string, req struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Dir          string           `json:"dir"`
	Config       []*ProcessConfig `json:"config"`
	Command      string           `json:"command"`
	Priority     int              `json:"priority"`
	StartRetries int              `json:"start_retries"`
	StopWaitSecs int              `json:"stop_wait_secs"`
	Groups       []string         `json:"groups"`
	Node         string           `json:"node"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	gs := map[string]struct{}{}
	for _, g := range req.Groups {
		if _, ok := processMgr.Groups[g]; !ok {
			// 创建路径
			ss := strings.Split(g, "/")
			var nav string
			for i := 0; i <= len(ss); i++ {
				nav = strings.Join(ss[0:i], "/")
				if _, ok := processMgr.Groups[nav]; !ok {
					processMgr.Groups[nav] = struct{}{}
				}
			}
		}
		gs[g] = struct{}{}
	}

	p, ok := processMgr.Process[req.ID]
	if !ok || p.State.Status != processStopped {
		done.result.Message = "当前状态不允许修改"
		return
	}

	p.Name = req.Name
	p.Dir = req.Dir
	p.Config = req.Config
	p.Command = req.Command
	p.Priority = req.Priority
	p.StartRetries = req.StartRetries
	p.StopWaitSecs = req.StartRetries
	p.Groups = gs
	p.Node = req.Node

	saveStore(snProcessMgr)
}

func (*processHandler) Delete(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok || p.State.Status != processStopped {
		done.result.Message = "当前状态不允许修改"
		return
	}

	delete(processMgr.Process, req.ID)
	saveStore(snProcessMgr)
}

func processTick() {
	rpcReq := map[string]*protocol.ProcessIsAliveReq{}
	for _, p := range processMgr.Process {
		if !(p.State.Status == processStopped ||
			p.State.Status == processExited) {
			req, ok := rpcReq[p.Node]
			if !ok {
				req = &protocol.ProcessIsAliveReq{
					Pid: map[int32]int32{},
				}
				rpcReq[p.Node] = req
			}
			req.Pid[int32(p.ID)] = p.State.Pid
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
			rpcResp := i.(*protocol.ProcessIsAliveResp)
			for id, b := range rpcResp.GetAlive() {
				p, ok := processMgr.Process[int(id)]
				if !ok {
					continue
				}
				if p.State.Status == processStarting {
					if b {
						p.State.Status = processRunning
					} else {
						p.State.Status = processExited
					}
				} else if p.State.Status == processRunning {
					if !b {
						p.State.Status = processExited
					}
				} else if p.State.Status == processStopping {
					if !b {
						p.State.Status = processStopped
					}
				}
			}
			saveStore(snProcessMgr)
		})
	}

}

func (*processHandler) Start(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok || !(p.State.Status == processStopped ||
		p.State.Status == processExited) {
		done.result.Message = "当前状态不允许启动"
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
		return
	}

	p.State = ProcessState{
		Status:    processStarting,
		Timestamp: NowUnix(),
	}

	configs := make(map[string]string, len(p.Config))
	for _, cfg := range p.Config {
		name := path.Join(p.Dir, strconv.Itoa(p.ID), cfg.Name)
		configs[name] = cfg.Context
	}

	cmd := strings.ReplaceAll(p.Command, "{{id}}", strconv.Itoa(p.ID))
	cmds := strings.Split(cmd, " ")
	rpcReq := &protocol.ProcessExecReq{
		Dir:    p.Dir,
		Name:   cmds[0],
		Args:   cmds[1:],
		Config: configs,
		Stderr: path.Join(p.Dir, strconv.Itoa(p.ID), "stderr.log"),
	}

	if err := center.Go(node, rpcReq, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
		if e != nil {
			done.result.Message = e.Error()
			done.Done()
			return
		}
		rpcResp := i.(*protocol.ProcessExecResp)
		if rpcResp.GetCode() != "" {
			done.result.Message = rpcResp.GetCode()
			p.State.Status = processStopped
		} else {
			p.State.Status = processRunning
			p.State.Pid = rpcResp.GetPid()
		}
		saveStore(snProcessMgr)
		done.Done()
	}); err != nil {
		log.Println(err)
		done.result.Message = err.Error()
		done.Done()
	} else {
		saveStore(snProcessMgr)
	}
}

func (*processHandler) Stop(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	p, ok := processMgr.Process[req.ID]
	if !ok || !(p.State.Status == processStopped ||
		p.State.Status == processStopping ||
		p.State.Status == processExited) {
		done.result.Message = "当前状态不允许停止"
		return
	}

	node, ok := nodes[p.Node]
	if !ok || !node.Online() {
		done.result.Message = "节点无服务"
		return
	}

	p.State.Status = processStopping

	rpcReq := &protocol.ProcessSignalReq{
		Pid:    p.State.Pid,
		Signal: int32(syscall.SIGTERM),
	}

	if err := center.Go(node, rpcReq, drpc.DefaultRPCTimeout, func(i interface{}, e error) {
		if e != nil {
			done.result.Message = e.Error()
			done.Done()
			return
		}
		rpcResp := i.(*protocol.ProcessSignalResp)
		if rpcResp.GetCode() != "" {
			done.result.Message = rpcResp.GetCode()
		}
		saveStore(snProcessMgr)
		done.Done()
	}); err != nil {
		log.Println(err)
		done.result.Message = err.Error()
		done.Done()
	} else {
		saveStore(snProcessMgr)
	}
}
