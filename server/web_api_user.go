package server

import (
	"encoding/json"
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

type userHandler struct {
}

func (*userHandler) Info(done *Done, user string) {
	log.Printf("%s by(%s)\n", done.route, user)
	defer func() { done.Done() }()
	//u, _ := getUser(user)
	//done.result.Data = struct {
	//	Name     string `json:"name"`
	//	Username string `json:"username"`
	//	Avatar   string `json:"avatar"`
	//}{
	//	Name:     u.Name,
	//	Username: u.Username,
	//	Avatar:   u.Avatar,
	//}
	str := `{
    "id": "4291d7da9005377ec9aec4a71ea837f", 
    "name": "管理员", 
    "username": "admin", 
    "password": "", 
    "avatar": "/avatar2.jpg", 
    "status": 1, 
    "telephone": "", 
    "lastLoginIp": "27.154.74.117", 
    "lastLoginTime": 1534837621348, 
    "creatorId": "admin", 
    "createTime": 1497160610259, 
    "merchantCode": "TLif2btpzg079h15bk", 
    "deleted": 0, 
    "roleId": "admin", 
    "role": {
        "id": "admin", 
        "name": "管理员", 
        "describe": "拥有所有权限", 
        "status": 1, 
        "creatorId": "system", 
        "createTime": 1497160610259, 
        "deleted": 0, 
        "permissions": [
            {
                "roleId": "admin", 
                "permissionId": "comment", 
                "permissionName": "评论管理", 
                "actions": "[]", 
                "actionEntitySet": [ ], 
                "actionList": [ ], 
                "dataAccess": null
            }
        ]
    }
}
`
	var info map[string]interface{}
	if err := json.Unmarshal([]byte(str), &info); err != nil {
		done.result.Code = 1
		done.result.Message = err.Error()
		return
	}
	done.result.Data = info
}

func (*userHandler) Nav(done *Done, user string) {
	log.Printf("%s by(%s)\n", done.route, user)
	defer func() { done.Done() }()

	done.result.Data = allNav
}

func (*userHandler) List(done *Done, user string, req struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}) {
	log.Printf("%s by(%s) %v\n", done.route, user, req)
	defer func() { done.Done() }()
	//if user != admin.Username {
	//	done.result.Code = 1
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
	//	done.result.Code = 1
	//	done.result.Message = "无权限"
	//	return
	//}

	if _, ok := userMgr.UserMap[req.Username]; ok {
		done.result.Code = 1
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
	//	done.result.Code = 1
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
