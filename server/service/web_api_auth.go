package service

import (
	"github.com/kataras/iris/v12"
	"log"
)

type Auth struct{}

func (*Auth) Login(req struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) (ret Result) {

	log.Printf("login %v\n", req)

	u, ok := getUser(req.Username)
	if !ok || u.Password != req.Password {
		ret.Code = 1
		ret.Message = "用户或密码错误"
		return
	}

	token := addToken(req.Username)
	ret.Data = struct {
		Token string `json:"token"`
	}{Token: token}
	return
}

/*
func (*Auth) Logout(user string) (ret Result) {
	log.Printf("logout %v\n", user)
	if _, ok := getUserTkn(user); ok {
		rmUserTkn(user)
	}
	return
}
*/
func (*Auth) Logout(ctx iris.Context) (ret Result) {
	tkn := ctx.GetHeader("Access-Token")
	if tkn == "" {
		return
	}

	if username, ok := getTknUser(tkn); ok {
		log.Printf("logout %s \n", username)
		rmTknUser(tkn)
	}
	return
}
