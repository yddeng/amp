package web

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"log"
)

const nav = `[
   {
    "name": "control",
    "parentId": 0,
    "id": 20,
    "meta": {
      "icon": "form",
      "title": "服务器管理"
    },
    "redirect": "/control/startpage",
    "component": "RouteView"
  },
  {
    "name": "startpage",
    "parentId": 20,
    "id": 21,
    "meta": {
      "title": "启动"
    },
    "component": "Startpage",
    "path": "/control/startpage"
  }
  ]`

func (*Handler) GetNav(ctx iris.Context) {
	username, ret := CheckToken(ctx)
	if ret != nil {
		ctx.JSON(ret)
		return
	}

	log.Printf("getNav user:%s \n", username)

	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(nav), &data); err != nil {
		ctx.Problem(NewProblem(iris.StatusInternalServerError, "", err.Error()))
		return
	}

	if _, err := ctx.JSON(Result{
		Data: data,
	}); err != nil {
		ctx.Problem(NewProblem(iris.StatusInternalServerError, "", err.Error()))
	}

}
