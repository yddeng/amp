package server

import (
	"fmt"
	"github.com/yddeng/dnet/drpc"
	"initial-server/protocol"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Cmd struct {
	Name     string              `json:"name"`
	Dir      string              `json:"dir"`
	Context  string              `json:"context"`
	Args     map[string]string   `json:"args"`
	User     string              `json:"user"`
	CreateAt int64               `json:"create_at"`
	CallNo   int                 `json:"call_no"`
	doing    map[string]struct{} // 节点正在执行
}

type CmdLog struct {
	ID       int    `json:"id"`
	CreateAt int64  `json:"create_at"` // 执行时间
	User     string `json:"user"`      // 执行用户
	Node     string `json:"node"`      // 执行的节点
	Context  string `json:"context"`   // 执行内容
	ResultAt int64  `json:"result_at"` // 执行结果时间
	Result   string `json:"result"`    // 执行结果
}

type CmdMgr struct {
	CmdMap  map[string]*Cmd      `json:"cmd_map"`
	CmdLogs map[string][]*CmdLog `json:"cmd_logs"`
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
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Cmd, 0, len(cmdMgr.CmdMap))
	for _, v := range cmdMgr.CmdMap {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].CreateAt > s[j].CreateAt
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = struct {
		PageNo     int    `json:"pageNo"`
		PageSize   int    `json:"pageSize"`
		TotalCount int    `json:"totalCount"`
		CmdList    []*Cmd `json:"cmdList"`
	}{PageNo: req.PageNo,
		PageSize:   req.PageSize,
		TotalCount: len(s),
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

	if _, ok := cmdMgr.CmdMap[req.Name]; ok {
		done.result.Code = 1
		done.result.Message = "名字重复"
		return
	}

	if len(cmdContextReg(req.Context)) != len(req.Args) {
		done.result.Code = 1
		done.result.Message = "变量与默认值数量不一致"
		return
	}

	cmd := &Cmd{
		Name:     req.Name,
		Dir:      req.Dir,
		Context:  req.Context,
		Args:     req.Args,
		User:     user,
		CreateAt: NowUnix(),
	}

	cmdMgr.CmdMap[req.Name] = cmd
	saveStore(snCmdMgr)
}

func (*cmdHandler) Delete(done *Done, user string, req struct {
	Name string `json:"name"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := cmdMgr.CmdMap[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令名"
		return
	}
	delete(cmdMgr.CmdMap, req.Name)
	delete(cmdMgr.CmdLogs, req.Name)
	saveStore(snCmdMgr)
}

func (*cmdHandler) Update(done *Done, user string, req struct {
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Context string            `json:"context"`
	Args    map[string]string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if cmd, ok := cmdMgr.CmdMap[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令名"
	} else {
		if len(cmdContextReg(req.Context)) != len(req.Args) {
			done.result.Code = 1
			done.result.Message = "变量与默认值数量不一致"
			return
		}

		cmd.Dir = req.Dir
		cmd.Context = req.Context
		cmd.Args = req.Args
		saveStore(snCmdMgr)
	}
}

const (
	cmdDefaultTimeout = 60
	cmdMinTimeout     = 10
	cmdMaxTimeout     = 86400
	cmdLogCapacity    = 20
)

func (*cmdHandler) Exec(done *Done, user string, req struct {
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Args    map[string]string `json:"args"`
	Node    string            `json:"node"`
	Timeout int               `json:"timeout"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	cmd, ok := cmdMgr.CmdMap[req.Name]
	if !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令"
		done.Done()
		return
	}

	node, ok := nodes[req.Node]
	if !ok || !node.Online() {
		done.result.Code = 1
		done.result.Message = "执行客户端不存在或不在线"
		done.Done()
		return
	}

	if cmd.doing == nil {
		cmd.doing = map[string]struct{}{}
	}
	if _, ok := cmd.doing[req.Node]; ok {
		done.result.Code = 1
		done.result.Message = "当前命令正在该节点上执行"
		done.Done()
		return
	}

	context := cmd.Context
	for k, v := range req.Args {
		context = strings.ReplaceAll(context, fmt.Sprintf("{{%s}}", k), v)
	}

	if len(cmdContextReg(context)) > 0 {
		done.result.Code = 1
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
	cmd.CallNo++
	callNo := cmd.CallNo
	cmdLog := &CmdLog{
		ID:       callNo,
		CreateAt: NowUnix(),
		User:     user,
		Node:     req.Node,
		Context:  context,
	}

	cmdMgr.CmdLogs[req.Name] = append(cmdMgr.CmdLogs[req.Name], cmdLog)
	if len(cmdMgr.CmdLogs[req.Name]) > cmdLogCapacity {
		cmdMgr.CmdLogs[req.Name] = cmdMgr.CmdLogs[req.Name][1:]
	}
	saveStore(snCmdMgr)

	cmdResult := func(cmdLog *CmdLog, ret string) {
		cmdLog.ResultAt = NowUnix()
		cmdLog.Result = ret
		if cmd.CallNo-cmdLog.ID < cmdLogCapacity {
			saveStore(snCmdMgr)
		}
	}

	cmd.doing[req.Node] = struct{}{}
	rpcReq := &protocol.CmdExecReq{
		Dir:     req.Dir,
		Name:    "/bin/sh",
		Args:    []string{"-c", context},
		Timeout: int32(req.Timeout),
	}
	timeout := time.Second*time.Duration(req.Timeout) + drpc.DefaultRPCTimeout
	if err := center.Go(node, rpcReq, timeout, func(i interface{}, e error) {
		rpcResp := i.(*protocol.CmdExecResp)
		if rpcResp.GetCode() != "" {
			done.result.Code = 1
			done.result.Message = rpcResp.GetCode()
			cmdResult(cmdLog, rpcResp.GetCode())
		} else {
			cmdResult(cmdLog, rpcResp.GetOutStr())
			done.result.Data = cmdLog
		}
		delete(cmd.doing, req.Node)
		done.Done()
	}); err != nil {
		log.Println(err)
		delete(cmd.doing, req.Node)
		cmdResult(cmdLog, err.Error())
		done.Done()
	}

}

func (*cmdHandler) Log(done *Done, user string, req struct {
	Name     string `json:"name"`
	PageNo   int    `json:"pageNo"`
	PageSize int    `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if logs, ok := cmdMgr.CmdLogs[req.Name]; !ok {
		done.result.Code = 1
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
