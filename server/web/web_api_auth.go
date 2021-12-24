package web

import (
	"github.com/kataras/iris/v12"
	"initialthree/node/node_gmmgr/webservice/problem"
	"log"
)

func login(context iris.Context) {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Login
	if err := context.ReadJSON(&req); err != nil {
		context.Problem(problem.New(iris.StatusBadRequest, "", err.Error()))
		return
	}

	ret := Result{}
	defer func() {
		if _, err := context.JSON(ret); err != nil {
			context.Problem(problem.New(iris.StatusInternalServerError, "", err.Error()))
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
