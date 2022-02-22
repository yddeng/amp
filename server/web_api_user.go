package server

import (
	"log"
	"sort"
)

type User struct {
	Username string              `json:"username,omitempty"`
	Password string              `json:"password,omitempty"`
	Routes   map[string]struct{} `json:"routes"`
}

type UserMgr struct {
	Admin   *User            `json:"admin"`
	UserMap map[string]*User `json:"user_map"`
}

type userHandler struct{}

func (*userHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	//log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	//if user != admin.Username {
	//	done.result.Message = "无权限"
	//	return
	//}

	s := make([]*User, 0, len(userMgr.UserMap))
	for _, v := range userMgr.UserMap {
		s = append(s, v)
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Username < s[j].Username
	})

	start, end := listRange(req.PageNo, req.PageSize, len(s))
	done.result.Data = struct {
		PageNo     int     `json:"pageNo"`
		PageSize   int     `json:"pageSize"`
		TotalCount int     `json:"totalCount"`
		UserList   []*User `json:"userList"`
	}{PageNo: req.PageNo,
		PageSize:   req.PageSize,
		TotalCount: len(s),
		UserList:   s[start:end]}
	return
}

func (*userHandler) Add(done *Done, user string, req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	//if user != admin.Username {
	//	done.result.Message = "无权限"
	//	return
	//}

	if _, ok := userMgr.UserMap[req.Username]; ok {
		done.result.Message = "用户名已存在"
		return
	}

	userMgr.UserMap[req.Username] = &User{
		Username: req.Username,
		Password: req.Password,
	}
	saveStore(snUserMgr)
}

func (*userHandler) Delete(done *Done, user string, req struct {
	Username []string `json:"username"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	//if user != admin.Username {
	//	done.result.Message = "无权限"
	//	return
	//}

	if len(req.Username) > 0 {
		for _, username := range req.Username {
			delete(userMgr.UserMap, username)
		}
		saveStore(snUserMgr)
	}
}
