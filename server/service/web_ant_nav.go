package service

import (
	"initial-sever/util"
)

type Nav struct {
	Name      string  `json:"name,omitempty"`
	Path      string  `json:"path,omitempty"`
	Id        int     `json:"id"`
	ParentId  int     `json:"parentId"`
	Meta      NavMeta `json:"meta,omitempty"`
	Redirect  string  `json:"redirect,omitempty"`
	Component string  `json:"component,omitempty"`
}

type NavMeta struct {
	Title        string `json:"title,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Show         bool   `json:"show,omitempty"`
	HideHeader   bool   `json:"hideHeader,omitempty"`
	HideChildren bool   `json:"hideChildren,omitempty"`
}

func findNav(routes map[string]struct{}) (navs []*Nav) {
	navs = make([]*Nav, 0, len(routes))
	for _, v := range allNav {
		if _, ok := routes[v.Name]; ok {
			navs = append(navs, v)
		}
	}
	return
}

func loadNav(filename string) error {
	return util.DecodeJsonFromFile(&allNav, filename)
}

var (
	allNav []*Nav
)

/*
	项目管理
		｜-- 启动模版管理
		｜-- 集群管理
				｜-- 集群列表
				｜-- 。。。 集群ID
*/

func newProjectNav(clus []*Cluster) []*Nav {
	navs := make([]*Nav, 0, len(clus)*4+5)
	navs = append(navs, &Nav{
		Name:     "project",
		Id:       10000,
		ParentId: 0,
		Meta: NavMeta{
			Title: "Initial项目",
			Icon:  "from",
		},
		Redirect:  "/project/cluster/cluster_list",
		Component: "Project",
	}, &Nav{
		Name:     "template",
		Id:       11000,
		ParentId: 10000,
		Path:     "/project/template",
		Meta: NavMeta{
			Title: "配置模版",
		},
		Component: "Template",
	}, &Nav{
		Name:     "cluster",
		Id:       12000,
		ParentId: 10000,
		Meta: NavMeta{
			Title: "集群组",
			Icon:  "from",
		},
		Path:      "/project/cluster",
		Component: "Cluster",
	})
	/*, &Nav{
		Name:     "cluster_list",
		Id:       12100,
		ParentId: 12000,
		Path:     "/project/cluster/cluster_list",
		Meta: NavMeta{
			Title: "集群列表",
			Icon:  "from",
		},
		Component: "ClusterList",
	})
	*/
	/*
		for i, v := range clus {
			n := i * 10
			id := fmt.Sprintf("%d", v.ID)
			navs = append(navs, &Nav{
				Name:     "cluster_group_" + id,
				Id:       12100 + n,
				ParentId: 12100,
				Meta: NavMeta{
					Title: "基本信息",
					Icon:  "from",
				},
				Redirect:  "/project/cluster/cluster_group_" + id,
				Component: "ClusterGroup",
			}, &Nav{
				Name:     "cluster_group_info_" + id,
				Id:       12100 + n + 1,
				ParentId: 12100 + n,
				Path:     "/project/cluster/cluster_group_" + id + "/cluster_group_info_" + id,
				Meta: NavMeta{
					Title: "基本信息",
				},
				Component: "ClusterGroupInfo",
			}, &Nav{
				Name:     "cluster_group_db_" + id,
				Id:       12100 + n + 2,
				ParentId: 12100 + n,
				Path:     "/project/cluster/cluster_group_" + id + "/cluster_group_db_" + id,
				Meta: NavMeta{
					Title: "数据库",
				},
				Component: "ClusterGroupDB",
			}, &Nav{
				Name:     "cluster_group_service_" + id,
				Id:       12100 + n + 3,
				ParentId: 12100 + n,
				Path:     "/project/cluster/cluster_group_" + id + "/cluster_group_service_" + id,
				Meta: NavMeta{
					Title: "服务管理",
				},
				Component: "ClusterGroupService",
			})
		}
	*/
	return navs
}
