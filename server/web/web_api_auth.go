package web

import (
	"initial-sever/task"
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

	rets := task.Wait(addToken, req.Username)
	token := rets[0].(string)

	ret.Data = struct {
		Token string `json:"token"`
	}{Token: token}
	return
}
