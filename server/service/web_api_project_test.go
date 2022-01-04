package service

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestTemplateHandler_Create(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()
	{
		//create
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/template/create", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Name string `json:"name"`
			Data string `json:"data"`
		}{Name: "center"})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		// update
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/template/update", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Name string `json:"name"`
			Data string `json:"data"`
		}{Name: "center"})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	//{
	//	//delete
	//	req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/template/delete", address), "POST")
	//	req.SetHeader("Access-Token", token)
	//	req.WriteJSON(struct {
	//		Name string `json:"name"`
	//	}{Name: "center"})
	//
	//	ret, err := req.ToString()
	//	t.Log(err, ret)
	//}

	{
		//list
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/template/list", address), "GET")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			PageNo   int `json:"pageNo"`
			PageSize int `json:"pageSize"`
		}{PageNo: 1, PageSize: 10})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

}

func TestClusterHandler_Create(t *testing.T) {
	startWebListener(t)

	ret := authLogin(t, "admin", "123456")
	t.Log(ret, gjson.Get(ret, "data.token").String())

	token := gjson.Get(ret, "data.token").String()

	{
		//create
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/template/create", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Name string `json:"name"`
			Data string `json:"data"`
		}{Name: "cluster", Data: `[Common]
centerAddr     	= [{0}]
dbConfig      	= ["flyfish@dir@{1}","flyfish@login@{1}","flyfish@game@{1}","flyfish@global@{1}", "flyfish@conflictzone@{1}","flyfish@nodelock@{1}"]
serverGroups    = [{2}]  # 服务的群组
cfgPathRoot     = "./configs"
excelPath       = "Excel"
questCfgPath    = "Quest"
wordsFilterPath     = "WordsFilter/wordsFilter.txt"
logicTimeSysCfgPath = "ScriptableObjects/TimeConfig.asset"
weatherSysCfgPath   = "ScriptableObjects/WeatherConfig.asset"`})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

	{
		//create
		req, _ := dhttp.NewRequest(fmt.Sprintf("http://%s/project/cluster/create", address), "POST")
		req.SetHeader("Access-Token", token)
		req.WriteJSON(struct {
			Desc    string   `json:"desc"`
			CfgTemp string   `json:"cfg_temp"`
			Args    []string `json:"args"`
		}{Desc: "cluster1", CfgTemp: "cluster", Args: []string{`"localhost:8010","localhost:8011"`, "localhost:10012", "1"}})

		ret, err := req.ToString()
		t.Log(err, ret)
	}

}
