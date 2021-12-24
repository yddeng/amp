package web

import (
	"github.com/kataras/iris/v12"
	"log"
)

func (*Handler) Login(ctx iris.Context) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Login
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.Problem(NewProblem(iris.StatusBadRequest, "", err.Error()))
		return
	}

	ret := Result{}
	defer func() {
		if _, err := ctx.JSON(ret); err != nil {
			ctx.Problem(NewProblem(iris.StatusInternalServerError, "", err.Error()))
		}
	}()

	log.Printf("login %v\n", req)
	if req.Password != "123456" {
		ret.Code = 1
		ret.Message = "密码错误"
		return
	}

	token := addToken(req.Username)
	ret.Data = struct {
		Token string `json:"token"`
	}{Token: token}
}
