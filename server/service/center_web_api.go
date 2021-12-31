package service

type nodeInfo struct {
	Name      string `json:"name"`
	Inet      string `json:"inet"`
	Net       string `json:"net"`
	Timestamp int64  `json:"timestamp"` // 登陆时间
	Online    bool   `json:"online"`
}

func getNodeInfo(cb func(nodes []nodeInfo)) {
	centerTaskQueue.Submit(func() {
		ninfos := make([]nodeInfo, 0, len(nodes))
		for _, n := range nodes {
			ninfos = append(ninfos, nodeInfo{
				Name:      n.Name,
				Inet:      n.Inet,
				Net:       n.Net,
				Timestamp: n.Timestamp,
				Online:    n.Online(),
			})
		}
		webTransQueue.Submit(cb, ninfos)
	})
}
