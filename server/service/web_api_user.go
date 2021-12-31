package service

import (
	"encoding/json"
	"log"
	"sort"
)

var (
	admin     *User
	userMap   = map[string]*User{}
	userSlice []*User
)

func sortUser() {
	sort.Slice(userSlice, func(i, j int) bool {
		return userSlice[i].Username < userSlice[j].Username
	})
}

func getUser(username string) (u *User, ok bool) {
	if admin.Username == username {
		return admin, true
	}
	u, ok = userMap[username]
	return
}

func addUser(username, password string) {
	u := &User{Username: username, Password: password}
	userMap[u.Username] = u
	userSlice = append(userSlice, u)
	sortUser()
	saveStore(snUser)
}

func deleteUser(username string) {
	if _, ok := userMap[username]; ok {
		delete(userMap, username)
		for i, u := range userSlice {
			if u.Username == username {
				userSlice = append(userSlice[:i], userSlice[i+1:]...)
			}
		}
	}
	saveStore(snUser)
}

type User struct {
	Username string              `json:"username,omitempty"`
	Password string              `json:"password,omitempty"`
	Routes   map[string]struct{} `json:"routes"`
}

type userHandler struct {
}

func (*userHandler) Info(done *Done, user string) {
	log.Printf("user/info by(%s) \n", user)
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
	log.Printf("user/nav by(%s) \n", user)
	defer func() { done.Done() }()

	done.result.Data = append(allNav, newProjectNav(nil)...)
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

	if userSlice == nil {
		userSlice = make([]*User, 0, len(userMap))
		for _, u := range userMap {
			if u.Username != admin.Username {
				userSlice = append(userSlice, u)
			}
		}
		sortUser()
	}

	start, end := listRange(req.PageNo, req.PageSize, len(userSlice))
	done.result.Data = struct {
		PageNo     int     `json:"pageNo"`
		PageSize   int     `json:"pageSize"`
		TotalCount int     `json:"totalCount"`
		UserList   []*User `json:"userList"`
	}{PageNo: req.PageNo,
		PageSize:   req.PageSize,
		TotalCount: len(userSlice),
		UserList:   userSlice[start:end]}
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

	if _, ok := getUser(req.Username); ok {
		done.result.Code = 1
		done.result.Message = "用户名已存在"
		return
	}

	addUser(req.Username, req.Password)
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

	for _, username := range req.Username {
		deleteUser(username)
	}
}
