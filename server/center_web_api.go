package server

import "time"

type nodeInfo struct {
	Name    string `json:"name"`
	Inet    string `json:"inet"`
	Net     string `json:"net"`
	LoginAt int64  `json:"login_at"` // 登陆时间
	Online  bool   `json:"online"`
}

func getNodeInfo(cb func(nodes []*nodeInfo)) {
	taskQueue.Submit(func() {
		ninfos := make([]*nodeInfo, 0, len(nodes))
		for _, n := range nodes {
			ninfos = append(ninfos, &nodeInfo{
				Name:    n.Name,
				Inet:    n.Inet,
				Net:     n.Net,
				LoginAt: n.LoginAt,
				Online:  n.Online(),
			})
		}
		time.Sleep(time.Second * 5)
		taskQueue.Submit(cb, ninfos)
	})
}
