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
	Name     string            `json:"name"`
	Dir      string            `json:"dir"`
	Context  string            `json:"context"`
	Args     map[string]string `json:"args"`
	User     string            `json:"user"`
	CreateAt int64             `json:"create_at"`
	doing    bool
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
	if !ok || cmd.doing {
		done.result.Code = 1
		done.result.Message = "不存在的命令或正在执行"
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

	cmd.doing = true
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
		} else {
			done.result.Data = struct {
				Output string `json:"output"`
			}{Output: rpcResp.GetOutStr()}
		}
		cmd.doing = false
		done.Done()
	}); err != nil {
		log.Println(err)
		cmd.doing = false
		done.Done()
	}

}
