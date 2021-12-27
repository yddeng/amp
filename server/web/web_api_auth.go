package web

import (
	"log"
)

type Auth struct{}

func (*Auth) Login(req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) (ret Result) {

	log.Printf("login %v\n", req)

	u, ok := getUser(req.Username)
	if !ok {
		ret.Code = 1
		ret.Message = "用户不存在"
		return
	}

	if u.Password != req.Password {
		ret.Code = 1
		ret.Message = "密码错误"
		return
	}

	token := addToken(req.Username)
	ret.Data = struct {
		Token string `json:"token"`
	}{Token: token}
	return
}

func (*Auth) Logout(user string) (ret Result) {
	log.Printf("logout %v\n", user)
	if _, ok := getUserTkn(user); ok {
		rmUserTkn(user)
	}
	return
}
