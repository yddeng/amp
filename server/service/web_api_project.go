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

type Template struct {
	Name     string `json:"name"`
	Shell    string `json:"shell"`
	Data     string `json:"data"`
	User     string `json:"user"`
	CreateAt int64  `json:"create_at"`
}

type templateHandler struct {
}

func (*templateHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Template, 0, len(temps))
	for _, v := range temps {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].CreateAt > s[j].CreateAt
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = s[start:end]
}

func (*templateHandler) Create(done *Done, user string, req struct {
	Name  string `json:"name"`
	Shell string `json:"shell"`
	Data  string `json:"data"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	if _, ok := temps[req.Name]; ok {
		done.result.Code = 1
		done.result.Message = "名字重复"
		return
	}

	temp := &Template{
		Name:     req.Name,
		Shell:    req.Shell,
		Data:     req.Data,
		User:     user,
		CreateAt: NowUnix(),
	}
	temps[req.Name] = temp
	saveStore(snTemplate)
}

func (*templateHandler) Update(done *Done, user string, req struct {
	Name  string `json:"name"`
	Shell string `json:"shell"`
	Data  string `json:"data"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	if temp, ok := temps[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的模版名"
	} else {
		temp.Shell = req.Shell
		temp.Data = req.Data
		temp.User = user
		saveStore(snTemplate)
	}
}

func (*templateHandler) Delete(done *Done, user string, req struct {
	Name string `json:"name"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	if _, ok := temps[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的模版名"
	}
	delete(temps, req.Name)
	saveStore(snTemplate)
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
