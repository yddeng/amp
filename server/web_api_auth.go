package server

import (
	"log"
)

type authHandler struct {
}

func (*authHandler) Login(done *Done, user string, req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) {

	log.Printf("%s %v\n", done.route, req)
	defer func() { done.Done() }()

	var u *User
	var ok bool
	if userMgr.Admin.Username == req.Username {
		u = userMgr.Admin
		ok = true
	} else {
		u, ok = userMgr.UserMap[req.Username]
	}

	if !ok || u.Password != req.Password {
		done.result.Message = "用户或密码错误"
		return
	}

	token := addToken(req.Username)
	done.result.Data = struct {
		Token string `json:"token"`
	}{Token: token}
	return
}

func (*authHandler) Logout(done *Done, user string) {
	log.Printf("%s by(%s) \n", done.route, user)
	defer func() { done.Done() }()
	if user == "" {
		return
	}
	rmUserTkn(user)
}
