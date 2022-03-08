package server

import (
	"amp/back-go/common"
	"amp/back-go/protocol"
	"github.com/yddeng/dnet/drpc"
	"log"
	"strings"
	"syscall"
)

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
	StartSecs      int64            `json:"start_secs"`
	StopWaitSecs   int64            `json:"stop_wait_secs"`
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
	p.StartSecs = req.StartSecs
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
	StartSecs      int64            `json:"start_secs"`
	StopWaitSecs   int64            `json:"stop_wait_secs"`
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
	p.StartSecs = req.StartSecs
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
