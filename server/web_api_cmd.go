package server

import (
	"fmt"
	"github.com/yddeng/dnet/drpc"
	"initial-server/protocol"
	"log"
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
	LogID    int                 `json:"log_id"`
	Logs     []*CmdLog           `json:"logs"`
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

type cmdHandler struct {
}

func (*cmdHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Cmd, 0, len(cmdMap))
	for _, v := range cmdMap {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].CreateAt > s[j].CreateAt
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = s[start:end]
}

func (*cmdHandler) Create(done *Done, user string, req struct {
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Context string            `json:"context"`
	Args    map[string]string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := cmdMap[req.Name]; ok {
		done.result.Code = 1
		done.result.Message = "名字重复"
		return
	}

	// todo 用正则表达式判断
	if strings.Count(req.Context, "{{") != len(req.Args) {
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

	cmdMap[req.Name] = cmd
	saveStore(snCmd)
}

func (*cmdHandler) Delete(done *Done, user string, req struct {
	Name string `json:"name"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := cmdMap[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令名"
		return
	}
	delete(cmdMap, req.Name)
	saveStore(snCmd)
}

func (*cmdHandler) Update(done *Done, user string, req struct {
	Name    string            `json:"name"`
	Dir     string            `json:"dir"`
	Context string            `json:"context"`
	Args    map[string]string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if cmd, ok := cmdMap[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令名"
	} else {
		// todo 用正则表达式判断
		if strings.Count(req.Context, "{{") != len(req.Args) {
			done.result.Code = 1
			done.result.Message = "变量与默认值数量不一致"
			return
		}

		cmd.Dir = req.Dir
		cmd.Context = req.Context
		cmd.Args = req.Args
		saveStore(snCmd)
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

	cmd, ok := cmdMap[req.Name]
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

	// todo 用正则表达式判断
	if strings.Contains(context, "{{") {
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
	cmd.LogID++
	logID := cmd.LogID
	cmdLog := &CmdLog{
		ID:       logID,
		CreateAt: NowUnix(),
		User:     user,
		Node:     req.Node,
		Context:  context,
	}
	cmd.Logs = append(cmd.Logs, cmdLog)
	if len(cmd.Logs) > cmdLogCapacity {
		cmd.Logs = cmd.Logs[1:]
	}
	saveStore(snCmd)

	cmdResult := func(logID int, ret string) {
		var _cmdLog *CmdLog
		for _, v := range cmd.Logs {
			if v.ID == logID {
				_cmdLog = v
			}
		}
		if _cmdLog != nil {
			_cmdLog.ResultAt = NowUnix()
			_cmdLog.Result = ret
			saveStore(snCmd)
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
			cmdResult(logID, rpcResp.GetCode())
		} else {
			done.result.Data = struct {
				Output string `json:"output"`
			}{Output: rpcResp.GetOutStr()}
			cmdResult(logID, rpcResp.GetOutStr())
		}
		delete(cmd.doing, req.Node)
		done.Done()
	}); err != nil {
		log.Println(err)
		delete(cmd.doing, req.Node)
		cmdResult(logID, err.Error())
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

	if cmd, ok := cmdMap[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的命令名"
	} else {
		start, end := listRange(req.PageNo, req.PageSize, len(cmd.Logs))
		done.result.Data = cmd.Logs[start:end]
	}
}
