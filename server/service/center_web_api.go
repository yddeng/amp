package service

type nodeInfo struct {
	Name    string `json:"name"`
	Inet    string `json:"inet"`
	Net     string `json:"net"`
	LoginAt int64  `json:"login_at"` // 登陆时间
	Online  bool   `json:"online"`
}

func getNodeInfo(cb func(nodes []*nodeInfo)) {
	centerTaskQueue.Submit(func() {
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
		webTransQueue.Submit(cb, ninfos)
	})
}

func getItemList(cb func(items []*Item)) {
	centerTaskQueue.Submit(func() {
		rets := make([]*Item, 0, len(itemMgr.Items))
		for _, v := range itemMgr.Items {
			rets = append(rets, &Item{
				ID:       v.ID,
				User:     v.User,
				Node:     v.Node,
				CreateAt: v.CreateAt,
				UpdateAt: v.UpdateAt,
				Online:   v.Online,
			})
		}
		webTransQueue.Submit(cb, rets)
	})
}
