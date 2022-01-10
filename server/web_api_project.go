package server

import (
	"fmt"
	"initial-server/util"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
)

type Cluster struct {
	ID       int    `json:"id"`
	Desc     string `json:"desc"`
	User     string `json:"user"`
	CreateAt int64  `json:"create_at"`
	Config   string `json:"config"`
}

type ClusterMgr struct {
	GenID    int              `json:"gen_id"`
	Clusters map[int]*Cluster `json:"clusters"`
}

type clusterHandler struct {
}

func (*clusterHandler) storeConfig(cluster *Cluster) {
	dir := fmt.Sprintf("cluster_%d", cluster.ID)
	file := "common.toml"
	curPath := path.Join(dataPath, dir)
	_ = os.MkdirAll(curPath, os.ModePerm)
	_ = ioutil.WriteFile(path.Join(curPath, file), []byte(cluster.Config), os.ModePerm)

}

func (*clusterHandler) deleteConfig(id int) {
	dir := fmt.Sprintf("cluster_%d", id)
	_ = os.RemoveAll(path.Join(dataPath, dir))
}

func (*clusterHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Cluster, 0, len(cluMgr.Clusters))
	for _, v := range cluMgr.Clusters {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID < s[j].ID
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = s[start:end]

}

func (this *clusterHandler) Create(done *Done, user string, req struct {
	Desc    string   `json:"desc"`
	CfgTemp string   `json:"cfg_temp"`
	Args    []string `json:"args"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	temp, ok := temps[req.CfgTemp]
	if !ok {
		done.result.Code = 1
		done.result.Message = "不存在的模版名"
		return
	}

	data := util.Replace(temp.Data, req.Args)
	cluMgr.GenID++
	cluster := &Cluster{
		ID:       cluMgr.GenID,
		Desc:     req.Desc,
		User:     user,
		Config:   data,
		CreateAt: NowUnix(),
	}
	cluMgr.Clusters[cluMgr.GenID] = cluster
	saveStore(snCluMgr)
	this.storeConfig(cluster)
}

func (this *clusterHandler) Delete(done *Done, user string, req struct {
	ID int `json:"id"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	if _, ok := cluMgr.Clusters[req.ID]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的集群"
	}
	delete(cluMgr.Clusters, req.ID)
	saveStore(snCluMgr)
	this.deleteConfig(req.ID)
}

type Template struct {
	Name     string `json:"name"`
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
	done.result.Data = struct {
		PageNo       int         `json:"pageNo"`
		PageSize     int         `json:"pageSize"`
		TotalCount   int         `json:"totalCount"`
		TemplateList []*Template `json:"templateList"`
	}{PageNo: req.PageNo,
		PageSize:     req.PageSize,
		TotalCount:   len(s),
		TemplateList: s[start:end],
	}
}

func (*templateHandler) Create(done *Done, user string, req struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
		Data:     req.Data,
		User:     user,
		CreateAt: NowUnix(),
	}
	temps[req.Name] = temp
	saveStore(snTemplate)
}

func (*templateHandler) Update(done *Done, user string, req struct {
	Name string `json:"name"`
	Data string `json:"data"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	if temp, ok := temps[req.Name]; !ok {
		done.result.Code = 1
		done.result.Message = "不存在的模版名"
	} else {
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
		return
	}
	delete(temps, req.Name)
	saveStore(snTemplate)
}

type ItemMgr struct {
	GenID int           `json:"gen_id"`
	Items map[int]*Item `json:"items"`
}

type Item struct {
	ID        int    `json:"id"`
	Desc      string `json:"desc"`
	User      string `json:"user"`
	CreateAt  int64  `json:"create_at"`
	UpdateAt  int64  `json:"update_at"`
	Online    bool   `json:"online"`
	Node      string `json:"node"`
	ClusterID int    `json:"cluster_id"`
	Shell     string `json:"shell"`
	Config    string `json:"config"`
}

type itemHandler struct {
}

func (*itemHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()

	s := make([]*Item, 0, len(itemMgr.Items))
	for _, v := range itemMgr.Items {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].UpdateAt > s[j].UpdateAt
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = s[start:end]
}
func (*itemHandler) Create(done *Done, user string, req struct {
	Desc      string `json:"desc"`
	Node      string `json:"node"`
	ClusterID int    `json:"cluster_id"`
	Shell     string `json:"shell"`
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
