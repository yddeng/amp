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
	if req.Username == "admin" && req.Password != "123456" {
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
