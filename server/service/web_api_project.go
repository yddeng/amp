package service

import (
	"log"
	"sort"
)

type Cluster struct {
	ID       int    `json:"id"`
	Desc     string `json:"desc"`
	User     string `json:"user"`
	Node     string `json:"node"`
	CreateAt int64  `json:"create_at"`
	CfgTemp  string `json:"cfg_temp"`
}

type ClusterMgr struct {
	GenID    int              `json:"gen_id"`
	Clusters map[int]*Cluster `json:"clusters"`
}

type clusterHandler struct {
}

func (*clusterHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*clusterHandler) Create(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*clusterHandler) Delete(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

type templateHandler struct {
}

func (*templateHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*templateHandler) Create(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*templateHandler) Update(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*templateHandler) Delete(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

type itemHandler struct {
}

func (*itemHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)

	getItemList(func(items []*Item) {
		sort.Slice(items, func(i, j int) bool {
			return items[i].UpdateAt > items[j].UpdateAt
		})
		start, end := listRange(req.PageNo, req.PageSize, len(items))
		done.result.Data = items[start:end]
		done.Done()
	})
}
func (*itemHandler) Create(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*itemHandler) Delete(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*itemHandler) Start(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}

func (*itemHandler) Single(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	done.Done()
}
