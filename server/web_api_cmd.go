package server

import (
	"amp/protocol"
	"fmt"
	"github.com/yddeng/dnet/drpc"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Cmd struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Dir      string              `json:"dir"`
	Context  string              `json:"context"`
	Args     map[string]string   `json:"args"`
	User     string              `json:"user"`
	UpdateAt int64               `json:"update_at"`
	CreateAt int64               `json:"create_at"`
	CallNo   int                 `json:"call_no"`
	doing    map[string]struct{} // 节点正在执行
}

type CmdMgr struct {
	Success int               `json:"success"`
	Failed  int               `json:"failed"`
	GenID   int               `json:"gen_id"`
	CmdMap  map[int]*Cmd      `json:"cmd_map"`
	CmdLogs map[int][]*CmdLog `json:"cmd_logs"`
}

// 以字母下划线开头，后接数字下划线和字母
func cmdContextReg(str string) map[string]struct{} {
	reg := regexp.MustCompile(`\{\{(_*[a-zA-Z]+[_a-zA-Z0-9]*)\}\}`)
	n := reg.FindAllString(str, -1)
	names := map[string]struct{}{}
	for _, name := range n {
		if _, ok := names[name]; !ok {
			names[name] = struct{}{}
		}
	}
	return names
}

type cmdHandler struct {
}

func (*cmdHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	//log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Cmd, 0, len(cmdMgr.CmdMap))
	for _, v := range cmdMgr.CmdMap {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].CallNo > s[j].CallNo
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = struct {
		PageNo     int    `json:"pageNo"`
		PageSize   int    `json:"pageSize"`
		TotalCount int    `json:"totalCount"`
		Success    int    `json:"success"`
		Failed     int    `json:"failed"`
		CmdList    []*Cmd `json:"cmdList"`
	}{PageNo: req.PageNo,
		PageSize:   req.PageSize,
		TotalCount: len(s),
		Success:    cmdMgr.Success,
		Failed:     cmdMgr.Failed,
		CmdList:    s[start:end],
	}
}

func (*cmdHandler) Create(done *Done, user string, req struct {
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Context string            `json:"context"`
	Args    map[string]string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	for _, cmd := range cmdMgr.CmdMap {
		if cmd.Name == req.Name {
			done.result.Message = "名字重复"
			return
		}
	}

	if len(cmdContextReg(req.Context)) != len(req.Args) {
		done.result.Message = "变量与默认值数量不一致"
		return
	}

	nowUnix := NowUnix()
	cmdMgr.GenID++
	id := cmdMgr.GenID
	cmd := &Cmd{
		ID:       id,
		Name:     req.Name,
		Dir:      req.Dir,
		Context:  req.Context,
		Args:     req.Args,
		User:     user,
		UpdateAt: nowUnix,
		CreateAt: nowUnix,
	}

	cmdMgr.CmdMap[id] = cmd
	saveStore(snCmdMgr)
}

func (*cmdHandler) Delete(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := cmdMgr.CmdMap[req.ID]; !ok {
		done.result.Message = "不存在的命令"
		return
	}
	delete(cmdMgr.CmdMap, req.ID)
	delete(cmdMgr.CmdLogs, req.ID)
	saveStore(snCmdMgr)
}

func (*cmdHandler) Update(done *Done, user string, req struct {
	ID      int               `json:"id"`
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Context string            `json:"context"`
	Args    map[string]string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	for _, cmd := range cmdMgr.CmdMap {
		if cmd.Name == req.Name && cmd.ID != req.ID {
			done.result.Message = "命令名重复"
			return
		}
	}

	if cmd, ok := cmdMgr.CmdMap[req.ID]; !ok {
		done.result.Message = "不存在的命令"
	} else {
		if len(cmdContextReg(req.Context)) != len(req.Args) {
			done.result.Message = "变量与默认值数量不一致"
			return
		}

		cmd.Name = req.Name
		cmd.Dir = req.Dir
		cmd.Context = req.Context
		cmd.Args = req.Args
		cmd.User = user
		cmd.UpdateAt = NowUnix()
		saveStore(snCmdMgr)
	}
}

const (
	cmdDefaultTimeout = 60
	cmdMinTimeout     = 10
	cmdMaxTimeout     = 86400
)

var cmdLogCapacity int = 10

type CmdLog struct {
	ID       int    `json:"id"`
	CreateAt int64  `json:"create_at"` // 执行时间
	User     string `json:"user"`      // 执行用户
	Dir      string `json:"dir"`       // 执行目录
	Node     string `json:"node"`      // 执行的节点
	Timeout  int    `json:"timeout"`   // 执行超时时间
	Context  string `json:"context"`   // 执行内容
	ResultAt int64  `json:"result_at"` // 执行结果时间
	Result   string `json:"result"`    // 执行结果
}

func (*cmdHandler) Exec(done *Done, user string, req struct {
	ID      int               `json:"id"`
	Dir     string            `json:"dir"`
	Args    map[string]string `json:"args"`
	Node    string            `json:"node"`
	Timeout int               `json:"timeout"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	cmd, ok := cmdMgr.CmdMap[req.ID]
	if !ok {
		done.result.Message = "不存在的命令"
		done.Done()
		return
	}

	node, ok := nodes[req.Node]
	if !ok || !node.Online() {
		done.result.Message = "执行客户端不存在或不在线"
		done.Done()
		return
	}

	if cmd.doing == nil {
		cmd.doing = map[string]struct{}{}
	}
	if _, ok := cmd.doing[req.Node]; ok {
		done.result.Message = "当前命令正在该节点上执行"
		done.Done()
		return
	}

	context := cmd.Context
	for k, v := range req.Args {
		context = strings.ReplaceAll(context, fmt.Sprintf("{{%s}}", k), v)
	}

	if len(cmdContextReg(context)) > 0 {
		done.result.Message = "命令中存在未赋值变量"
		done.Done()
		return
	}

	// 超时时间
	if req.Timeout <= 0 {
		req.Timeout = cmdDefaultTimeout
	} else if req.Timeout < cmdMinTimeout {
		req.Timeout = cmdMinTimeout
	} else if req.Timeout > cmdMaxTimeout {
		req.Timeout = cmdMaxTimeout
	}

	// 执行日志

	cmdLog := &CmdLog{
		CreateAt: NowUnix(),
		Timeout:  req.Timeout,
		User:     user,
		Node:     req.Node,
		Dir:      req.Dir,
		Context:  context,
	}

	cmdResult := func(cmdLog *CmdLog, ok bool, ret string) {
		if ok {
			cmdMgr.Success += 1
		} else {
			cmdMgr.Failed += 1
		}
		cmdLog.ResultAt = NowUnix()
		cmdLog.Result = ret
		saveStore(snCmdMgr)
	}

	rpcReq := &protocol.CmdExecReq{
		Dir:     req.Dir,
		Name:    "/bin/sh",
		Args:    []string{"-c", context},
		Timeout: int32(req.Timeout),
	}
	timeout := time.Second*time.Duration(req.Timeout) + drpc.DefaultRPCTimeout
	if err := center.Go(node, rpcReq, timeout, func(i interface{}, e error) {
		if e != nil {
			done.result.Message = e.Error()
			cmdResult(cmdLog, false, e.Error())
			done.Done()
			return
		}
		rpcResp := i.(*protocol.CmdExecResp)
		if rpcResp.GetCode() != "" {
			done.result.Message = rpcResp.GetCode()
			cmdResult(cmdLog, false, rpcResp.GetCode())
		} else {
			cmdResult(cmdLog, true, rpcResp.GetOutStr())
			done.result.Data = cmdLog
		}
		delete(cmd.doing, req.Node)
		done.Done()
	}); err != nil {
		log.Println(err)
		done.result.Message = err.Error()
		done.Done()
	} else {
		cmd.doing[req.Node] = struct{}{}
		cmd.CallNo++
		cmdLog.ID = cmd.CallNo
		cmdMgr.CmdLogs[req.ID] = append([]*CmdLog{cmdLog}, cmdMgr.CmdLogs[req.ID]...)
		if len(cmdMgr.CmdLogs[req.ID]) > cmdLogCapacity {
			cmdMgr.CmdLogs[req.ID] = cmdMgr.CmdLogs[req.ID][:cmdLogCapacity]
		}
		saveStore(snCmdMgr)
	}

}

func (*cmdHandler) Log(done *Done, user string, req struct {
	ID       int `json:"id"`
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	//log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if logs, ok := cmdMgr.CmdLogs[req.ID]; !ok {
		done.result.Message = "不存在的命令名"
	} else {
		start, end := listRange(req.PageNo, req.PageSize, len(logs))
		done.result.Data = struct {
			PageNo     int       `json:"pageNo"`
			PageSize   int       `json:"pageSize"`
			TotalCount int       `json:"totalCount"`
			LogList    []*CmdLog `json:"logList"`
		}{PageNo: req.PageNo,
			PageSize:   req.PageSize,
			TotalCount: len(logs),
			LogList:    logs[start:end],
		}
	}
}
