package server

import (
	"amp/back-go/protocol"
	"log"
	"sort"
)

type nodeInfo struct {
	Name    string              `json:"name"`
	Inet    string              `json:"inet"`
	Net     string              `json:"net"`
	LoginAt int64               `json:"login_at"` // 登陆时间
	Online  bool                `json:"online"`
	State   *protocol.NodeState `json:"state"`
}

type nodeHandler struct{}

func (*nodeHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	//log.Printf("%s by(%s) %v\n", done.route, user, req)

	s := make([]*nodeInfo, 0, len(nodes))
	for _, n := range nodes {
		s = append(s, &nodeInfo{
			Name:    n.Name,
			Inet:    n.Inet,
			Net:     n.Net,
			LoginAt: n.LoginAt,
			Online:  n.Online(),
			State:   n.nodeState,
		})
	}
	sort.Slice(s, func(i, j int) bool {
		if s[i].Online == s[j].Online {
			return s[i].LoginAt > s[j].LoginAt
		} else {
			return s[i].Online
		}
	})

	start, end := listRange(req.PageNo, req.PageSize, len(nodes))
	done.result.Data = struct {
		PageNo     int         `json:"pageNo"`
		PageSize   int         `json:"pageSize"`
		TotalCount int         `json:"totalCount"`
		NodeList   []*nodeInfo `json:"nodeList"`
	}{PageNo: req.PageNo,
		PageSize:   req.PageSize,
		TotalCount: len(s),
		NodeList:   s[start:end],
	}
	done.Done()
}

func (*nodeHandler) Remove(done *Done, user string, req struct {
	Name string `json:"name"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	n, ok := nodes[req.Name]
	if !ok || n.Online() {
		done.result.Message = "当前状态不允许移除"
		return
	}

	delete(nodes, req.Name)
	saveStore(snNode)
}
